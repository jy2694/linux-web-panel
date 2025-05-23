import { auth } from '$lib/stores/auth';
import type { AuthState } from '$lib/stores/auth';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

export const load = async ({ url }: { url: URL }) => {
  if (browser) {
    // 인증 스토어 초기화
    auth.initialize();

    // 로그인 페이지가 아닌 경우에만 인증 체크
    if (url.pathname !== '/login') {
      let isAuthenticated = false;
      auth.subscribe((state: AuthState) => {
        isAuthenticated = state.isAuthenticated;
      })();
      if (!isAuthenticated) {
        goto('/login');
      }
    }
  }

  return {};
}; 