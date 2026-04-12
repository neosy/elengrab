let activeMenus = new Set();
const isPWA =
      window.matchMedia('(display-mode: standalone)').matches ||
      window.matchMedia('(display-mode: fullscreen)').matches ||
      window.navigator.standalone === true;

// Close all opened menus except the provided one
export function closeAllMenus(except = null) {
  document.querySelectorAll('.menu.show').forEach(menu => {
    if (menu !== except) {
      menu.classList.remove('show');
      menu._activeTrigger = null;
      activeMenus.delete(menu);
    }
  });
}

// Generic menu initializer
export function initMenu(config) {
  const {
    triggerSelector,
    menuId,
    actions = {},
    position,
    buildUrl,
    shouldOpen,
    beforeOpen
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

  document.body.addEventListener('click', (e) => {
    const trigger = e.target.closest(triggerSelector);
    if (!trigger) return;

    e.stopPropagation();

    const isSameTrigger = trigger === menu._activeTrigger;

    closeAllMenus(menu);

    if (isSameTrigger && menu.classList.contains('show')) {
      menu.classList.remove('show');
      menu._activeTrigger = null;
      activeMenus.delete(menu);
      return;
    }

    menu._activeTrigger = trigger;

    position?.(menu, trigger);

    if (buildUrl) {
      const url = buildUrl(menu, trigger);
      if (url) {
        menu.setAttribute('hx-get', url);
        window.htmx?.process(menu);
      }
    }

    if (shouldOpen?.(trigger) === false) return;

    beforeOpen?.(menu, trigger);

    menu.innerHTML = "";
    menu.classList.add('show');
    activeMenus.add(menu);

    window.htmx?.trigger(menu, 'manual');
  });

  /** Handle clicks inside the menu */
  menu.addEventListener('click', async (e) => {
    const item = e.target.closest('[data-action]');
    if (!item) return;

    e.stopPropagation();
    try {
      await handleMenuAction(item);
    } catch (err) {
      console.error('Menu action failed:', err);
    }

    menu.classList.remove('show');
    menu._activeTrigger = null;
    activeMenus.delete(menu);
  });

  /** Global click to close menu if clicked outside */
  document.addEventListener('click', (e) => {
    const target = e.target;
    if (menu.contains(target) || menu._activeTrigger?.contains(target)) return;
    menu.classList.remove('show');
    menu._activeTrigger = null;
    activeMenus.delete(menu);
  });
}

function kebabToCamel(str) {
  return str.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
}