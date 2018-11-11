const uidKey = 'userID';

export function getAuth() {
  return {
    nonce: window.localStorage.getItem('nonce'),
    userID: window.localStorage.getItem(uidKey),
  };
}

export function isLoggedIn() {
  let nonce = window.localStorage.getItem('nonce');
  let userID = window.localStorage.getItem(uidKey);

  return nonce !== null && userID !== null;
}

export function logIn() {
  window.location.href = '/api/auth/twitch/user?redirect='+window.location.pathname;
}

export function logOut() {
  // TODO: server-side log out as well
  window.localStorage.removeItem('nonce');
  window.localStorage.removeItem(uidKey);

  window.location.reload();
}

export function loadAuth() {
  let nonce = window.localStorage.getItem('nonce');
  let user_id = window.localStorage.getItem(uidKey);
  if (window.location.hash) {
    let hash = window.location.hash.substring(1);
    let hashParts = hash.split(';');
    let hashMap = {};
    for (let i = 0; i < hashParts.length; ++i) {
      let p = hashParts[i].split('=');
      if (p.length == 2) {
        hashMap[p[0]] = p[1];
      }
    }

    if ('nonce' in hashMap && 'user_id' in hashMap) {
      nonce = hashMap['nonce'];
      user_id = hashMap['user_id'];
      window.localStorage.setItem('nonce', nonce);
      window.localStorage.setItem(uidKey, user_id);
    }

    history.replaceState(
      '', document.title,
      window.location.pathname + window.location.search);
  }

  return { nonce: nonce, userId: user_id };

}
