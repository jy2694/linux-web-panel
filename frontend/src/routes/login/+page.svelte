<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let username = '';
  let password = '';
  let error = '';
  let loading = false;

  onMount(() => {
    // 이미 로그인되어 있다면 홈으로 리다이렉트
    if ($auth.isAuthenticated) {
      goto('/');
    }
  });

  async function handleSubmit() {
    loading = true;
    error = '';

    try {
      const success = await auth.login(username, password);
      if (success) {
        goto('/');
      } else {
        error = '로그인에 실패했습니다.';
      }
    } catch (e) {
      error = '로그인 중 오류가 발생했습니다.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        Linux Web Panel
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        로그인하여 계속하세요
      </p>
    </div>
    <form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
      <div class="rounded-md shadow-sm -space-y-px">
        <div>
          <label for="username" class="sr-only">사용자 이름</label>
          <input
            id="username"
            name="username"
            type="text"
            required
            bind:value={username}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="사용자 이름"
          />
        </div>
        <div>
          <label for="password" class="sr-only">비밀번호</label>
          <input
            id="password"
            name="password"
            type="password"
            required
            bind:value={password}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="비밀번호"
          />
        </div>
      </div>

      {#if error}
        <div class="text-red-500 text-sm text-center">{error}</div>
      {/if}

      <div>
        <button
          type="submit"
          disabled={loading}
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {#if loading}
            로그인 중...
          {:else}
            로그인
          {/if}
        </button>
      </div>
    </form>
  </div>
</div> 