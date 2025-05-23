<script lang="ts">
  import { onMount } from 'svelte';
  import axios from 'axios';

  interface User {
    id: string;
    username: string;
    role: string;
    created_at: string;
    updated_at: string;
  }

  let users: User[] = [];
  let error = '';
  let loading = false;
  let showAddModal = false;
  let newUser = {
    username: '',
    password: '',
    role: 'user'
  };

  async function fetchUsers() {
    loading = true;
    try {
      const response = await axios.get('http://localhost:8080/api/users');
      users = response.data;
    } catch (err) {
      error = '사용자 목록을 가져오는데 실패했습니다.';
      console.error(err);
    } finally {
      loading = false;
    }
  }

  async function addUser() {
    try {
      await axios.post('http://localhost:8080/api/users', newUser);
      showAddModal = false;
      newUser = { username: '', password: '', role: 'user' };
      await fetchUsers();
    } catch (err) {
      error = '사용자 추가에 실패했습니다.';
      console.error(err);
    }
  }

  async function deleteUser(id: string) {
    if (!confirm('정말로 이 사용자를 삭제하시겠습니까?')) return;
    try {
      await axios.delete(`http://localhost:8080/api/users/${id}`);
      await fetchUsers();
    } catch (err) {
      error = '사용자 삭제에 실패했습니다.';
      console.error(err);
    }
  }

  onMount(fetchUsers);
</script>

<div class="container mx-auto px-4 py-8">
  <div class="flex justify-between items-center mb-8">
    <h1 class="text-3xl font-bold">사용자 관리</h1>
    <button
      class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors duration-150"
      on:click={() => showAddModal = true}
    >
      사용자 추가
    </button>
  </div>

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
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">사용자명</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">역할</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">생성일</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">작업</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#if loading}
            <tr>
              <td colspan="4" class="px-6 py-4 text-center text-gray-500">로딩 중...</td>
            </tr>
          {:else if users.length === 0}
            <tr>
              <td colspan="4" class="px-6 py-4 text-center text-gray-500">사용자가 없습니다.</td>
            </tr>
          {:else}
            {#each users as user}
              <tr>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{user.username}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {user.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800'}">
                    {user.role}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(user.created_at).toLocaleString()}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <button
                    class="text-red-600 hover:text-red-900 transition-colors duration-150"
                    on:click={() => deleteUser(user.id)}
                  >
                    삭제
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

{#if showAddModal}
  <div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
      <h2 class="text-xl font-bold mb-6">사용자 추가</h2>
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">사용자명</label>
          <input
            type="text"
            bind:value={newUser.username}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="사용자명을 입력하세요"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">비밀번호</label>
          <input
            type="password"
            bind:value={newUser.password}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="비밀번호를 입력하세요"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">역할</label>
          <select
            bind:value={newUser.role}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="user">일반 사용자</option>
            <option value="admin">관리자</option>
          </select>
        </div>
      </div>
      <div class="mt-6 flex justify-end space-x-3">
        <button
          class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors duration-150"
          on:click={() => showAddModal = false}
        >
          취소
        </button>
        <button
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors duration-150"
          on:click={addUser}
        >
          추가
        </button>
      </div>
    </div>
  </div>
{/if} 