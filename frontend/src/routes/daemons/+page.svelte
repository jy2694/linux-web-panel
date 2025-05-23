<script lang="ts">
  import { onMount } from 'svelte';
  import axios from 'axios';

  interface Daemon {
    name: string;
    status: string;
    description: string;
  }

  let daemons: Daemon[] = [];
  let error = '';
  let loading = false;

  async function fetchDaemons() {
    loading = true;
    try {
      const response = await axios.get('http://localhost:8080/api/daemons');
      daemons = response.data;
    } catch (err) {
      error = '데몬 목록을 가져오는데 실패했습니다.';
      console.error(err);
    } finally {
      loading = false;
    }
  }

  async function controlDaemon(name: string, action: 'start' | 'stop' | 'restart') {
    try {
      await axios.post(`http://localhost:8080/api/daemons/${name}/${action}`);
      await fetchDaemons();
    } catch (err) {
      error = `${name} 데몬 ${action}에 실패했습니다.`;
      console.error(err);
    }
  }

  onMount(fetchDaemons);
</script>

<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-8">데몬 관리</h1>

  {#if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
      {error}
    </div>
  {/if}

  <div class="card">
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">이름</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">상태</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">설명</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">작업</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#if loading}
            <tr>
              <td colspan="4" class="px-6 py-4 text-center text-gray-500">로딩 중...</td>
            </tr>
          {:else if daemons.length === 0}
            <tr>
              <td colspan="4" class="px-6 py-4 text-center text-gray-500">데몬이 없습니다.</td>
            </tr>
          {:else}
            {#each daemons as daemon}
              <tr>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{daemon.name}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {daemon.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                    {daemon.status}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-gray-500">{daemon.description}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <div class="flex space-x-3">
                    <button
                      class="bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600 transition-colors duration-150"
                      on:click={() => controlDaemon(daemon.name, 'start')}
                    >
                      시작
                    </button>
                    <button
                      class="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600 transition-colors duration-150"
                      on:click={() => controlDaemon(daemon.name, 'stop')}
                    >
                      중지
                    </button>
                    <button
                      class="bg-yellow-500 text-white px-3 py-1 rounded hover:bg-yellow-600 transition-colors duration-150"
                      on:click={() => controlDaemon(daemon.name, 'restart')}
                    >
                      재시작
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

<style>
  :global(body) {
    background-color: #f3f4f6;
  }
</style> 