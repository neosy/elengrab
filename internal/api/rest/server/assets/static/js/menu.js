import { isPWA } from "./browser.js";

const MENU_HASH = "#menu"

const CLASS_NAMES = {
  menuOverlayShow: 'menu-overlay--show',
  menuMobileLayout: 'menu--mobile-layout',
  menuAction: 'menu-action',

  menuShow: 'menu--show',
  menuUp: 'menu--up',
  menuDown: 'menu--down',
  menuScroll: 'menu--scroll',

  getMenuStateClasses() {
    return [
      this.menuShow,
      this.menuUp,
      this.menuDown,
      this.menuScroll,
    ];
  },
}

const CLASS_SELECTORS = Object.fromEntries(
    Object.entries(CLASS_NAMES)
        .filter(([, value]) => typeof value === "string")
        .map(([key, value]) => [key, `.${value}`])
);

const DOM_IDS = {
  menuOverlay: 'menu-overlay',
}

const DOM_ELEMENTS = {
  menuOverlay: null,
}

function initDomElements() {
  DOM_ELEMENTS.menuOverlay = document.getElementById(DOM_IDS.menuOverlay);
}

let activeMenus = new Set();
const menuState = new Map();

// Close all opened menus except the provided one
export function closeAllMenus(except = null, fnAfterClose) {
  document.querySelectorAll(CLASS_SELECTORS.menuShow).forEach(menu => {
    if (menu !== except) {
      closeMenu(menu, fnAfterClose);
    }
  });
}

// Generic menu initializer
/**
 * Menu configuration
 *
 * @typedef {(item: HTMLElement) => void} menuAction
 * @typedef {Object.<string, menuAction>} menuActions 
 *
 * Required:
 * @property {string} triggerSelector
 * @property {string} menuId
 * @property {menuActions} actions
 *
 * Hooks (DOM lifecycle):
 * @property {(menu: HTMLElement, trigger: HTMLElement) => void} [beforeOpen]
 * @property {(menu: HTMLElement) => void} [afterClose]
 *
 * Rendering:
 * @property {(el: HTMLElement, id: string) => void} [beforeRenderElement]
 * @property {(el: HTMLElement, trigger: HTMLElement) => void} [position]
 *
 * Logic:
 * @property {(trigger: HTMLElement) => string} [shouldOpen]
 * @property {(menu: HTMLElement, trigger: HTMLElement) => string} [buildUrl]
 * @property {() => boolean} [isMobile]
 * @property {boolean} [isOverlay]
 * @property {boolean} [isStatic] Whether the menu content is static and does not require server loading.
 */
export function initMenu(config) {
  const {
    triggerSelector,
    menuId,
    actions = {},
    position,
    isMobile = null,
    isOverlay = true,
    isStatic = false,
    beforeRenderElement,
    buildUrl,
    shouldOpen,
    beforeOpen,
    afterClose
  } = config;

  const menu = document.getElementById(menuId);
  if (!menu) return;

  async function handleMenuAction(item) {
    const action = kebabToCamel(item.dataset.action);
    const handler = actions[action] || actions.default || defaultNavigate;
    await handler(item, menu);
  }

  function defaultNavigate(item) {
    if (item.dataset.url) {
      const isNewTab = item.dataset.newTab === 'true';
      // For links, open in new tab if specified and not in PWA mode; otherwise navigate in same tab
      // In PWA mode, we generally want to stay in the app, so we ignore newTab for internal links
      if (isNewTab && !isPWA) {
        window.open(item.dataset.url, '_blank', 'noopener,noreferrer');
      } else {
        window.location.href = item.dataset.url;
      }
    }
  }

  function isMobileDefault() {
    return window.matchMedia('(max-width: 768px)').matches;
  }

  function applyMobile(menu) {
    const props = ['top', 'left', 'right', 'bottom', 'transform', 'position'];
    props.forEach(prop => {
      menu.style[prop] = '';
    });
    menu.classList.remove(CLASS_NAMES.menuMobileLayout);
    menu.classList.add(CLASS_NAMES.menuMobileLayout);
  }  

  function syncMenuWithHash() {
      if (location.hash === MENU_HASH) {
          return;
      }

      if (menuState.get(menu.id) === true) {
          closeMenu(menu, afterClose);
      }
  }

  document.body.addEventListener('click', (e) => {
    const trigger = e.target.closest(triggerSelector);
    if (!trigger) return;

    e.stopPropagation();

    const isSameTrigger = trigger === menu._activeTrigger;

    closeAllMenus(menu, afterClose);

    if (isSameTrigger && menu.classList.contains(CLASS_NAMES.menuShow)) {
      closeMenu(menu, afterClose);
      return;
    }

    menuState.set(menu.id, true);
    location.hash = MENU_HASH

    menu._activeTrigger = trigger;

    if (buildUrl) {
      const url = buildUrl(menu, trigger);
      if (url) {
        menu.setAttribute('hx-get', url);
        window.htmx?.process(menu);
      }
    }

    if (shouldOpen?.(trigger) === false) {
      closeMenu(menu, afterClose);
      return;
    }

    beforeOpen?.(menu, trigger);

    if (isOverlay) {
      DOM_ELEMENTS.menuOverlay && (DOM_ELEMENTS.menuOverlay.classList.add(CLASS_NAMES.menuOverlayShow));
    }
   
    if (!isStatic) {
      menu.innerHTML = "";
    }

    activeMenus.add(menu);
    window.htmx?.trigger(menu, 'manual');

    const handleMenuReady = () => {
      menu.querySelectorAll(CLASS_SELECTORS.menuAction).forEach(el => {
        const id = el.id.replace('menu-action-', '');
        beforeRenderElement?.(el, id);
      });

      if (isMobile ? isMobile() : isMobileDefault()) {
        applyMobile(menu);
      } else {
        position?.(menu, trigger);
      }

      requestAnimationFrame(() => {
        menu.classList.add(CLASS_NAMES.menuShow);
      });
    };

    if (isStatic) {
      handleMenuReady();
    } else {
      menu.addEventListener('htmx:afterSettle', handleMenuReady, { once: true });
    }
  });

  /** Handle clicks inside the menu */
  menu.addEventListener('click', async (e) => {
    if (!menuState.get(menu.id)) return;

    closeMenu(menu, afterClose);

    const item = e.target.closest('[data-action]');
    if (!item) return;

    e.stopPropagation();

    try {
      await handleMenuAction(item);
    } catch (err) {
      console.error('Menu action failed:', err);
    }
  });

  /** Global click to close menu if clicked outside */
  document.addEventListener('click', (e) => {
    if (!menuState.get(menu.id)) return;

    const target = e.target;
    if (menu.contains(target) || menu._activeTrigger?.contains(target)) return;
    closeMenu(menu, afterClose);
  });

  // Global ESC handler
  document.addEventListener("keydown", (event) => {
    if (!menuState.get(menu.id)) return;

    if (event.key === "Escape" && location.hash === MENU_HASH) {
      closeMenu(menu, afterClose);
    }
  });

  if (location.hash === MENU_HASH) {
      history.replaceState(null, "", location.pathname + location.search);
  }

  window.addEventListener('hashchange', syncMenuWithHash);
}

function closeMenu(menu, fnAfterClose) {
  menuState.set(menu.id, false);
  
  history.replaceState(null, "", location.pathname + location.search);

  menu.classList.remove(...CLASS_NAMES.getMenuStateClasses());

  if (DOM_ELEMENTS.menuOverlay) {
    DOM_ELEMENTS.menuOverlay.classList.remove(CLASS_NAMES.menuOverlayShow)
  }
  
  menu.style.removeProperty('max-height');
  
  menu._activeTrigger = null;

  activeMenus.delete(menu);
  fnAfterClose?.(menu);
}

function kebabToCamel(str) {
  if (!str) {return}
  return str.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
}


export function positionFloatingMenu(menu, btn) {
  if (!menu || !btn) return;

  const rect = btn.getBoundingClientRect();
  const gap = 6;

  menu.style.maxHeight = '500px';
  const menuHeight = menu.scrollHeight;
  menu.style.removeProperty('max-height');

  const spaceBelow = window.innerHeight - rect.bottom;
  const spaceAbove = rect.top - gap;

  if (spaceBelow >= menuHeight + gap) {
      // open down
      menu.style.top = `${rect.bottom + gap}px`;
      menu.style.bottom = 'auto';

      menu.classList.add(CLASS_NAMES.menuDown);
  } else if (spaceAbove >= menuHeight) {
      // open up
      menu.style.bottom = `${window.innerHeight - rect.top + gap}px`;
      menu.style.top = 'auto';

      menu.classList.add(CLASS_NAMES.menuUp);
  } else {
    // doesn't fit either side
    if (spaceBelow >= spaceAbove) {
        menu.style.top = `${rect.bottom + gap}px`;
        menu.style.bottom = 'auto';
        menu.style.maxHeight = `${spaceBelow}px`;

        menu.classList.add(CLASS_NAMES.menuDown);
    } else {
        menu.style.bottom = `${window.innerHeight - rect.top + gap}px`;
        menu.style.top = 'auto';
        menu.style.maxHeight = `${spaceAbove}px`;

        menu.classList.add(CLASS_NAMES.menuUp);
    }

    menu.classList.add(CLASS_NAMES.menuScroll);
  }

  menu.style.left = 'auto';
  menu.style.right = `${window.innerWidth - rect.right}px`;
}

document.addEventListener('DOMContentLoaded', () => {
  initDomElements();
})
