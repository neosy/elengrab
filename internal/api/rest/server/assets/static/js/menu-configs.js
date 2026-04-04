import { initMenu } from './menu.js';

// Account menu config
const accountMenuConfig = {
  triggerSelector: '#user-avatar',
  menuId: 'account-menu',

  actions: {
    default(item) {
      if (item.dataset.url) {
        window.location.href = item.dataset.url;
      }
    }
  },

  beforeOpen(menu, trigger) {
    if (trigger.dataset.action !== 'menu') return;
    menu.innerHTML = "";
  }
};

// Row menu config
const rowMenuConfig = {
  triggerSelector: '.result__menu-button',
  menuId: 'row-menu',

  actions: {
    async copyLink(item) {
      const res = await fetch(item.dataset.url, { method: 'POST' });
      if (!res.ok) throw new Error('Request failed');

      const data = await res.json();

      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(data.url);
        console.log('Link copied:', data.url);
      } else {
        console.log('Link not copied (Clipboard API not supported):', data.url);
        alert('Clipboard API not supported');
      }
    },

    async delete(item) {
      const res = await fetch(item.dataset.url, { method: 'DELETE' });
      if (!res.ok) throw new Error('Request failed');
    }
  },

  position(menu, btn) {
    const rect = btn.getBoundingClientRect();

    menu.style.top = `${rect.bottom + 6}px`;
    menu.style.left = 'auto';
    menu.style.right = `${window.innerWidth - rect.right}px`;
  },

  buildUrl(menu, btn) {
    let url = menu.dataset.menuUrl;
    const fileId = btn.dataset.fileId;

    if (!url || !fileId) return null;

    return url.replace('%7bfileId%7d', fileId);
  }
};

// Initialize all menus
export function initMenus() {
  initMenu(accountMenuConfig);
  initMenu(rowMenuConfig);
}