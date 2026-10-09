import { initMenu } from './menu.js';
import { initMenus as configInitMenus } from './pages-list.menu-configs.js';
import { DOM_IDS } from './index-page.dom.js';

export const DOM_ELEMENTS = {
  mediaURLInput: null,
}

function initDomElements() {
  DOM_ELEMENTS.mediaURLInput = document.getElementById(DOM_IDS.mediaURLInput);;
}

// Upload menu config
const uploadMenuConfig = {
  triggerSelector: '#header-actions-upload-menu-button',
  menuId: 'upload-menu',
  isStatic: true,

  actions: {
    uploadUrl(item) {
      if (!DOM_ELEMENTS.mediaURLInput) return;

      DOM_ELEMENTS.mediaURLInput.scrollIntoView({
          behavior: "smooth",
          block: "center",
      });

      DOM_ELEMENTS.mediaURLInput.focus();
    }
  },

  shouldOpen(trigger) {
    return true;
  },

  beforeOpen(menu, trigger) {
    return true;
  }
};

export function initMenus() {
  initDomElements();

  configInitMenus();
  initMenu(uploadMenuConfig);
}