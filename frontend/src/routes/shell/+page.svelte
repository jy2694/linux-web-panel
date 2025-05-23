<script lang="ts">
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';

  let ws: WebSocket;
  let connected = false;
  let error = '';
  let output = '';
  let input = '';

  function connect() {
    if (!browser) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/ws/shell`;
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
      connected = true;
      error = '';
    };
    
    ws.onclose = () => {
      connected = false;
      error = '연결이 끊어졌어요.';
    };
    
    ws.onerror = () => {
      error = '연결 오류가 발생했어요.';
    };
    
    ws.onmessage = (event) => {
      output += event.data;
    };
  }

  function sendCommand() {
    if (ws && ws.readyState === WebSocket.OPEN && input.trim()) {
      ws.send(input + '\n');
      input = '';
    }
  }

  onMount(() => {
    if (browser) {
      connect();
    }
  });
</script>

<div class="max-w-4xl mx-auto p-4">
  <div class="flex justify-between items-center mb-4">
    <h1 class="text-2xl font-bold">웹 터미널</h1>
    {#if !connected}
      <button
        class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        on:click={connect}
      >
        재연결
      </button>
    {/if}
  </div>

  {#if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {error}
    </div>
  {/if}

  <div class="bg-black text-white font-mono p-4 rounded-lg shadow-lg overflow-auto" style="height: 500px;">
    <pre>{output}</pre>
  </div>

  <div class="mt-4 flex">
    <input
      type="text"
      bind:value={input}
      on:keydown={(e) => e.key === 'Enter' && sendCommand()}
      class="flex-1 bg-gray-800 text-white px-4 py-2 rounded-l focus:outline-none"
      placeholder="명령어를 입력하세요..."
    />
    <button
      class="bg-blue-600 text-white px-4 py-2 rounded-r hover:bg-blue-700"
      on:click={sendCommand}
    >
      실행
    </button>
  </div>
</div> 