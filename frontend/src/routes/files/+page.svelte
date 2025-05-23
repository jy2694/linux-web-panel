<script lang="ts">
  import { onMount } from 'svelte';
  import axios from 'axios';

  interface FileItem {
    name: string;
    path: string;
    type: 'file' | 'directory';
    size: number;
    modified: string;
  }

  let currentPath = '/';
  let files: FileItem[] = [];
  let error = '';
  let loading = false;

  async function fetchFiles(path: string) {
    loading = true;
    try {
      const response = await axios.get(`http://localhost:8080/api/files?path=${encodeURIComponent(path)}`);
      files = response.data;
      currentPath = path;
    } catch (err) {
      error = '파일 목록을 가져오는데 실패했습니다.';
      console.error(err);
    } finally {
      loading = false;
    }
  }

  function formatSize(size: number): string {
    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
    if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(1)} MB`;
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`;
  }

  function formatDate(date: string): string {
    return new Date(date).toLocaleString();
  }

  function navigateUp() {
    const parentPath = currentPath.split('/').slice(0, -1).join('/') || '/';
    fetchFiles(parentPath);
  }

  onMount(() => fetchFiles('/'));
</script>

<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-8">파일 관리</h1>

  {#if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
      {error}
    </div>
  {/if}

  <div class="card">
    <div class="p-4 border-b border-gray-200">
      <div class="flex items-center space-x-2">
        <button
          class="text-blue-600 hover:text-blue-800 disabled:text-gray-400 disabled:cursor-not-allowed"
          on:click={navigateUp}
          disabled={currentPath === '/'}
        >
          상위 디렉토리
        </button>
        <span class="text-gray-500">/</span>
        <span class="font-mono text-gray-700">{currentPath}</span>
      </div>
    </div>

    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">이름</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">크기</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">수정일</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#if loading}
            <tr>
              <td colspan="3" class="px-6 py-4 text-center text-gray-500">로딩 중...</td>
            </tr>
          {:else if files.length === 0}
            <tr>
              <td colspan="3" class="px-6 py-4 text-center text-gray-500">파일이 없습니다.</td>
            </tr>
          {:else}
            {#each files as file}
              <tr 
                class="hover:bg-gray-50 cursor-pointer transition-colors duration-150"
                on:click={() => file.type === 'directory' && fetchFiles(file.path)}
              >
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <span class="text-gray-900">{file.name}</span>
                    {#if file.type === 'directory'}
                      <span class="ml-2 text-gray-400">/</span>
                    {/if}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {file.type === 'file' ? formatSize(file.size) : '-'}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(file.modified)}
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