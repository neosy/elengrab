import { shareLink } from './helper.js';
import { initMenu } from './menu.js';
import * as notify from './notifications.js';

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
  triggerSelector: '.media-result__menu-button',
  menuId: 'row-menu',

  actions: {
    async shareLink(item) {
      const result = await shareLink({
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
          //alert('Copy this link: ' + result.url);
          break;

        case 'error':
          notify.show("Failed to share link", notify.notifyType.ERROR);
          console.error('Share error:', result.error);
          //alert('Failed to share link');
          break;
      }
    },

    async copyLink(item) {
      const res = await fetch(item.dataset.url, { method: 'POST' });
      if (!res.ok) throw new Error('Request failed');

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
  },

  beforeOpen(menu, trigger) {
    document.documentElement.classList.add("ui-blocking-active");
    document.body.classList.add('ui-blocking-active');
    return true;
  },

  afterClose(menu) {
    document.documentElement.classList.remove("ui-blocking-active");
    document.body.classList.remove('ui-blocking-active');
    return true;
  }
};

// Initialize all menus
export function initIndexMenus() {
  initMenu(accountMenuConfig);
  initMenu(rowMenuConfig);
}