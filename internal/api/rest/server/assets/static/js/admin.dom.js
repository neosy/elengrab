export const DOM_IDS = {
    userRolesList: "userRolesList",
    
    rowUser: (id) => `row-user-${id}`,
}

export const DOM_ELEMENTS = {
    main: null,

    usersPanel: null,
    usersTable: null,
    userDetailPanel: null,
    userDetailUserName: null,
    userDetailBody: null,
    userDetailCloseBtn: null,
    userDetailCancelBtn: null,
    userDetailSaveBtn: null,
};

export function initDomElements() {
    DOM_ELEMENTS.main = document.getElementById("main");

    DOM_ELEMENTS.usersPanel = document.getElementById("usersPanel");
    DOM_ELEMENTS.usersTable = document.getElementById("usersTable");
    DOM_ELEMENTS.userDetailPanel = document.getElementById("userDetailPanel");
    DOM_ELEMENTS.userDetailUserName = document.getElementById("userDetailUserName");
    DOM_ELEMENTS.userDetailBody = document.getElementById("userDetailBody");
    DOM_ELEMENTS.userDetailCloseBtn = document.getElementById("userDetailCloseBtn");
    DOM_ELEMENTS.userDetailCancelBtn = document.getElementById("userDetailCancelBtn");
    DOM_ELEMENTS.userDetailSaveBtn = document.getElementById("userDetailSaveBtn");
};
