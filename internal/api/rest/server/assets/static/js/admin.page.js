import { DOM_ELEMENTS, initDomElements } from "./admin.dom.js";
import { DOM_IDS } from './constants.js';
import * as utils from './utils.js';
import * as browser from './browser.js';
import * as notify from './notifications.js';

const adminSection = {
    DASHBOARD: 'dashboard',
    USERS: 'users',
    SETTINGS: 'settings',
};

const EDIT_HASH = "#edit";
const PANEL_HASH = "#panel";

const ROUTES = [
    ['/admin/users', adminSection.USERS],
    ['/admin/settings', adminSection.SETTINGS],
];

function getCurrentAdminSection() {
    const path = window.location.pathname;

    if (path === '/admin') {
        return adminSection.DASHBOARD;
    }

    const route = ROUTES.find(([prefix]) => path.startsWith(prefix));

    return route?.[1] ?? null;
}

document.addEventListener('DOMContentLoaded', () => {
    initDomElements();

    // Initialize viewport height sync (fixes mobile PWA viewport issues)
    browser.initViewportHeightVar();

    // Initialize content panel
    initContentPanel();

    const currentSection = getCurrentAdminSection();

    switch (currentSection) {
        case adminSection.USERS:
            initUserDetailPanel();
            break;
    }
});

function initContentPanel() {
    const panel = document.querySelector('.content__panel');
    if (!panel) return;

    document.body.classList.add('content-open');
    syncHash(true);

    const backButton = panel.querySelector('.back-button--content');
    if (backButton) {
        backButton.addEventListener('click', () => {
            closeContentPanel();
        });
    }

    window.addEventListener('hashchange', () => syncHash());
}

function closeContentPanel() {
    document.body.classList.remove('content-open');
    syncHash();
}

function syncHash(isInit = false) {
    const isContentOpen = document.body.classList.contains('content-open');
    const isDetailViewOpen = document.body.classList.contains('detail-view-open');

    if (isInit) {
        location.hash = PANEL_HASH;
        return;
    }

    if (location.hash === "") {
        if (isContentOpen) {
            closeContentPanel();
            return;
        }
    }

    let hash = location.hash;
    if (hash === EDIT_HASH && !isDetailViewOpen) {
        hash = "";
    }

    if (isContentOpen) {
        if (hash === "") {
            hash = PANEL_HASH;
        }
    } else {
        hash = "";
    }

    if (hash === "") {
        history.replaceState(null, '', window.location.pathname + window.location.search);
    } else {
        location.hash = hash;
    }

    if (hash !== EDIT_HASH && isDetailViewOpen) {
        closeUserDetailPanel();
    }
}

function initUserDetailPanel() {
    const userDetailPanel = DOM_ELEMENTS.userDetailPanel;
    const userDetailBody = DOM_ELEMENTS.userDetailBody;
    const userDetailUserName = DOM_ELEMENTS.userDetailUserName;
    const usersTable = DOM_ELEMENTS.usersTable;

    if (!userDetailPanel || !userDetailBody || !usersTable) return;

    usersTable.addEventListener('click', (event) => {
        const button = event.target.closest('.user-edit-btn');
        if (!button) return;

        const { userId, userName } = button.dataset;

        openUserDetailPanel(userId, userName, userDetailPanel, userDetailUserName, userDetailBody);
    });    

    DOM_ELEMENTS.userDetailCloseBtn.addEventListener('click', closeUserDetailPanel);
    DOM_ELEMENTS.userDetailCancelBtn.addEventListener('click', closeUserDetailPanel);
    DOM_ELEMENTS.userDetailSaveBtn.addEventListener('click', saveUserDetail);
}

function openUserDetailPanel(userId, userName, panel, userDetailUserName, detailBody) {
    if (!panel || !userDetailUserName || !detailBody) {
        return;
    }

    closeUserDetailPanel();

    let path = detailBody.dataset.queryPath;
    path = path.replace('{userId}', userId);

    userDetailUserName.textContent = userName;
    detailBody.setAttribute('hx-get', path);

    htmx.process(panel);

    const onAfterRequest = (event) => {
        if (event.target !== detailBody) {
            return;
        }

        detailBody.removeEventListener('htmx:afterRequest', onAfterRequest);

        if (event.detail.successful) {
            panel.classList.add('active');
            document.body.classList.add('detail-view-open');
            location.hash = EDIT_HASH;
        } else {
            console.error('Failed to load user details', event.detail);
            notify.show(
                'Failed to load user details. Please try again later.',
                notify.notifyType.ERROR
            );
        }
    };

    detailBody.addEventListener('htmx:afterRequest', onAfterRequest);

    htmx.trigger(detailBody, 'loadUserDetail');
}

function closeUserDetailPanel() {
    if (DOM_ELEMENTS.userDetailPanel) {
        DOM_ELEMENTS.userDetailPanel.classList.remove('active');
    }
    document.body.classList.remove('detail-view-open');

    if (location.hash === EDIT_HASH) {
        history.back();
    }
}

async function saveUserDetail() {
    try {
        const rolesList = document.getElementById(DOM_IDS.userRolesList);
        if (!rolesList) return;

        const userId = rolesList.dataset.userId;

        const payload = {
            userId: userId,
            roles: getSelectedRoles(rolesList),
        };

        const response = await fetch(rolesList.dataset.updateRolesPath, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            let message = 'Failed to update user roles.';

            try {
                const error = await response.json();
                message = error.message || message;
            } catch (_) {
            }

            throw new Error(message);
        }

        closeUserDetailPanel();

        notify.show(
            'User roles updated successfully.',
            notify.notifyType.SUCCESS
        );

        const userRow = document.getElementById(DOM_IDS.rowUser(userId));
        if (userRow) {
            htmx.trigger(userRow, 'loadUserTableRow');
        }
    } catch (error) {
        notify.show(
            error.message || 'Failed to update user roles. Please try again later.',
            notify.notifyType.ERROR
        );
    }
}

function getSelectedRoles(rolesList) {
    return Array.from(
        rolesList.querySelectorAll('input[type="checkbox"]:checked')
    ).map(cb => cb.dataset.roleId);
}
