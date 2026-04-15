import { shareLink } from './helper.js';
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
  triggerSelector: '.result__menu-button',
  menuId: 'row-menu',

  actions: {
    async shareLink(item) {
      const result = await shareLink({
        endpoint: item.dataset.url,
        method: 'POST',
        title: "Elengrab",
        text: 'Watch this'
      });      

      switch (result.status) {
        case 'shared':
          console.log('Shared:', result.url);
          break;

        case 'copied':
          console.log('Link copied:', result.url);
          break;

        case 'manual':
          console.log('Manual copy needed:', result.url);
          alert('Copy this link: ' + result.url);
          break;

        case 'error':
          console.error('Share error:', result.error);
          alert('Failed to share link');
          break;
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
export function initIndexMenus() {
  initMenu(accountMenuConfig);
  initMenu(rowMenuConfig);
}