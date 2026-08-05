import { CLASS_NAMES } from './constants.js';
import * as share from './share.js';
import { initMenu, DOM_CLASSES as MENU_CLASSES, positionFloatingMenu} from './menu.js';
import * as notify from './notifications.js';
import * as view from './index.view.js';

// Account menu config
const accountMenuConfig = {
  triggerSelector: '#user-avatar',
  menuId: 'accountMenu',

  actions: {
    default(item) {
      if (item.dataset.url) {
        window.location.href = item.dataset.url;
      }
    }
  },

  shouldOpen(trigger) {
    if (trigger.dataset.action !== 'menu') return false;
    return true;
  },

  beforeOpen(menu, trigger) {
    menu.innerHTML = "";
    return true;
  }
};

// Settings menu config
const settingsMenuConfig = {
  triggerSelector: '#footerSettingsLink',
  menuId: 'settingsMenu',

  actions: {
    gridView(item) {
      const isGridView = view.toggleGridView();
    }
  },

  beforeRenderElement(el, id) {
    if (!el ) return;

    if (id === "grid-view") {
      const span = el.querySelector('.menu-action__content-text-value');
      let text = el.dataset.text

      if (span && text) {
        text = `${text}: ${view.getGridView() ? 'on' : 'off'}`;
        span.textContent = text;
      }
    }
  },

  position(menu, btn) {
    const rect = btn.getBoundingClientRect();

    menu.style.left = 'auto';
    menu.style.right = `${window.innerWidth - rect.right}px`;

    requestAnimationFrame(() => {
      menu.style.bottom = `${window.innerHeight - rect.top + 6}px`;
      menu.style.top = 'auto';      
    });
  },

  isMobile() {
    return false;
  },

  isOverlay: false,

  shouldOpen(trigger) {
    if (trigger.dataset.action !== 'menu') return false;
    return true;
  },

  beforeOpen(menu, trigger) {
    menu.innerHTML = "";
    return true;
  }
};

// Row menu config
const rowMenuConfig = {
  triggerSelector: '.media-result__menu-button',
  menuId: 'rowMenu',

  actions: {
    async errorInfo(item) {
      const errorText = item.dataset.text;
      if (errorText) {
        notify.show(errorText, notify.notifyType.ERROR, 7*1000);
      } else {
        notify.show('Error information is not available', notify.notifyType.ERROR);
      }
    },

    async shareLink(item) {
      const result = await share.shareLink({
        endpoint: item.dataset.url,
        method: 'POST',
        title: "Elengrab",
        text: ''
      });      

      switch (result.status) {
        case 'shared':
          notify.show("Link shared", notify.notifyType.SUCCESS);
          console.log('Shared:', result.url);
          break;

        case 'copied':
          notify.show("Link copied", notify.notifyType.SUCCESS);
          console.log('Link copied:', result.url);
          break;

        case 'manual':
          notify.show('Copy this link: ' + result.url, notify.notifyType.INFO, 8000);
          console.log('Copy this link:', result.url);
          break;

        case 'error':
          if (result.error !== "") {
            notify.show(result.error.message, notify.notifyType.ERROR);  
          } else {
            notify.show("Failed to share link", notify.notifyType.ERROR);  
          }
          console.error('Share error:', result.error);
          break;
      }
    },

    async createLink(item) {
      const res = await fetch(item.dataset.url, { method: 'POST' });
      if (!res.ok) {
        const data = await res.json();
        console.info(data);
        if (data && typeof data === "object" && "message" in data) {
            notify.show(data.message, notify.notifyType.ERROR);
        }
        throw new Error('Request failed');
      }

      const data = await res.json();

      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(data.url);
        notify.show("Link copied.", notify.notifyType.SUCCESS);
        console.log('Link copied:', data.url);
      } else {
        notify.show("Clipboard API not supported", notify.notifyType.ERROR);
        notify.show('Copy this link: ' + data.url, notify.notifyType.INFO, 8000);
        console.log('Link not copied (Clipboard API not supported):', data.url);
        //alert('Clipboard API not supported');
      }
    },

    async copyLink(item) {
      const res = await fetch(item.dataset.url, { method: 'GET' });
      if (!res.ok) {
        const data = await res.json();
        console.info(data);
        if (data && typeof data === "object" && "message" in data) {
            notify.show(data.message, notify.notifyType.ERROR);
        }
        throw new Error('Request failed');
      }

      const data = await res.json();

      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(data.url);
        notify.show("Link copied.", notify.notifyType.SUCCESS);
        console.log('Link copied:', data.url);
      } else {
        notify.show("Clipboard API not supported", notify.notifyType.ERROR);
        notify.show('Copy this link: ' + data.url, notify.notifyType.INFO, 8000);
        console.log('Link not copied (Clipboard API not supported):', data.url);
      }
    },

    async deleteLink(item) {
      const res = await fetch(item.dataset.url, { method: 'DELETE' });
      if (!res.ok) {
        const data = await res.json();
        console.info(data);
        if (data && typeof data === "object" && "message" in data) {
            notify.show(data.message, notify.notifyType.ERROR);
        }
        throw new Error('Request failed');
      }

      notify.show("Link deleted.", notify.notifyType.SUCCESS);
    },

    async refresh(item) {
      const res = await fetch(item.dataset.url, { method: 'POST' });
      if (!res.ok) {
        const data = await res.json();
        if (data && typeof data === "object" && "message" in data) {
            notify.show(data.message, notify.notifyType.ERROR);
        }
        notify.show('Request failed', notify.notifyType.ERROR);
        throw new Error('Request failed');
      }
    },

    async delete(item) {
      const res = await fetch(item.dataset.url, { method: 'DELETE' });
      if (!res.ok) {
        const data = await res.json();
        if (data && typeof data === "object" && "message" in data) {
            notify.show(data.message, notify.notifyType.ERROR);
        }
        notify.show('Request failed', notify.notifyType.ERROR);
        throw new Error('Request failed');
      }
    }
  },

  position(menu, btn) {
    positionFloatingMenu(menu, btn);
  },

  buildUrl(menu, btn) {
    let url = menu.dataset.menuUrl;
    const itemId = btn.dataset.itemId;

    if (!url || !itemId) return null;

    return url.replace('%7bitemId%7d', itemId);
  },

  beforeOpen(menu, trigger) {
    document.documentElement.classList.add(CLASS_NAMES.ui.blockingActive);
    document.body.classList.add(CLASS_NAMES.ui.blockingActive);
    return true;
  },

  afterClose(menu) {
    document.documentElement.classList.remove(CLASS_NAMES.ui.blockingActive);
    document.body.classList.remove(CLASS_NAMES.ui.blockingActive);
    return true;
  }
};

// Initialize all menus
export function initIndexMenus() {
  initMenu(accountMenuConfig);
  initMenu(settingsMenuConfig);
  initMenu(rowMenuConfig);
}