<script lang="ts">
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';
  import '../app.css';

  function handleLogout() {
    auth.logout();
    window.location.href = '/login';
  }

  $: isLoginPage = $page.url.pathname === '/login';
</script>

<div class="min-h-screen flex flex-col bg-gray-100">
  {#if !isLoginPage}
    <!-- 네비게이션 바 -->
    <nav class="bg-gray-800 sticky top-0 z-50">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <a href="/" class="text-white text-xl font-bold">Linux Web Panel</a>
          </div>
          <div class="flex items-center space-x-4">
            <a href="/overview" class="nav-link {$page.url.pathname === '/overview' ? 'active' : ''}">
              시스템 개요
            </a>
            <a href="/shell" class="nav-link {$page.url.pathname === '/shell' ? 'active' : ''}">
              웹 터미널
            </a>
            <a href="/daemons" class="nav-link {$page.url.pathname === '/daemons' ? 'active' : ''}">
              데몬 관리
            </a>
            <a href="/files" class="nav-link {$page.url.pathname === '/files' ? 'active' : ''}">
              파일 관리
            </a>
            <a href="/users" class="nav-link {$page.url.pathname === '/users' ? 'active' : ''}">
              사용자 관리
            </a>
            <button
              on:click={handleLogout}
              class="ml-4 px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              로그아웃
            </button>
          </div>
        </div>
      </div>
    </nav>
  {/if}

  <!-- 메인 컨텐츠 -->
  <main class="flex-grow py-6">
    <slot />
  </main>

  {#if !isLoginPage}
    <!-- 푸터 -->
    <footer class="bg-white border-t border-gray-200">
      <div class="container mx-auto px-4 py-6">
        <p class="text-center text-gray-600 text-sm">
          &copy; 2024 Linux Web Panel. All rights reserved.
        </p>
      </div>
    </footer>
  {/if}
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif;
  }
</style> 