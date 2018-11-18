import {readCookie} from "./cookie";

export function loggedInUsername() {
  return readCookie('pb2username') || '???';
}

export function isLoggedIn() {
  return readCookie('pb2sessionid') !== null;
}

export function logIn() {
  window.location.href = '/api/auth/twitch/user?redirect='+window.location.pathname;
}

export function logOut() {
  window.location.href = '/logout?redirect='+window.location.pathname;
}
