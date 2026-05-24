const isPWA =
      window.matchMedia('(display-mode: standalone)').matches ||
      window.matchMedia('(display-mode: fullscreen)').matches ||
      window.navigator.standalone === true;

const MENU_SHOW_CLASS = 'menu--show';
const MENU_OVERLAY_SHOW_CLASS = 'menu-overlay--show';
const MENU_MOBILE_LAYOUT_CLASS = 'menu--mobile-layout';
const MENU_ACTION_CLASS = 'menu-action';
const MENU_HASH = "#menu"

const menuOverlay = document.getElementById("menu-overlay");

let activeMenus = new Set();
const menuState = new Map();

// Close all opened menus except the provided one
export function closeAllMenus(except = null, fnAfterClose) {
  document.querySelectorAll('.'+MENU_SHOW_CLASS).forEach(menu => {
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
 */
export function initMenu(config) {
  const {
    triggerSelector,
    menuId,
    actions = {},
    position,
    isMobile,
    isOverlay = true,
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
    menu.classList.remove(MENU_MOBILE_LAYOUT_CLASS);
    menu.classList.add(MENU_MOBILE_LAYOUT_CLASS);
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

    if (isSameTrigger && menu.classList.contains(MENU_SHOW_CLASS)) {
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
      menuOverlay && (menuOverlay.classList.add(MENU_OVERLAY_SHOW_CLASS));
    }
   
    menu.innerHTML = "";
    activeMenus.add(menu);

    window.htmx?.trigger(menu, 'manual');

    menu.addEventListener(
      'htmx:afterSettle',
      () => {
        menu.querySelectorAll(`.${MENU_ACTION_CLASS}`).forEach(el => {
          const id = el.id.replace('menu-action-', '');
          beforeRenderElement?.(el, id);
        });

        menu.classList.add(MENU_SHOW_CLASS);

        if (isMobile ? isMobile() : isMobileDefault()) {
          applyMobile(menu);
        } else {
          position?.(menu, trigger);
        }
      },
      { once: true }
    );        
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

  menu.classList.remove(MENU_SHOW_CLASS);
  menuOverlay && (menuOverlay.classList.remove(MENU_OVERLAY_SHOW_CLASS));

  menu._activeTrigger = null;

  activeMenus.delete(menu);
  fnAfterClose?.(menu);
}

function kebabToCamel(str) {
  if (!str) {return}
  return str.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
}
