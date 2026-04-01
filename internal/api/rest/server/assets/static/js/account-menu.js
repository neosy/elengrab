export function initAccountMenu() {
  const avatar = document.getElementById('user-avatar');
  const menu = document.getElementById('account-menu');

  if (!avatar || !menu) return;

  const getAvatarAction = () => avatar.dataset.action; // 'menu', 'none', 'login'

  // Click on avatar
  avatar.addEventListener('click', (e) => {
    e.stopPropagation();

    const action = getAvatarAction();

    switch(action) {
      case 'menu':
        menu.classList.toggle('show');
        break;
      case 'login':
      case 'none':
      default:
        break;
    } 
  });

  // Close menu when clicking outside
  document.addEventListener('click', (e) => {
    if (!menu.contains(e.target) && !avatar.contains(e.target)) {
      menu.classList.remove('show');
    }
  });

  menu.addEventListener('click', (e) => {
    const item = e.target.closest('[data-href]');
    if (!item) return;

    e.stopPropagation();

    const url = item.dataset.href;
    if (url) {
      window.location.href = url;
    }
  });
}


