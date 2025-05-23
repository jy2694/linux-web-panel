<script lang="ts">
  import { onMount } from 'svelte';
  import axios from 'axios';

  let systemInfo = {
    hostname: '',
    uptime: '',
    os: '',
    kernel: '',
    cpu: { usage: 0 },
    memory: { total: 0, used: 0, free: 0 },
    disk: { total: 0, used: 0, free: 0 }
  };

  let cpuPercentage = 0;
  let memoryPercentage = 0;
  let diskPercentage = 0;

  async function fetchSystemInfo() {
    try {
      const response = await axios.get('http://localhost:8080/api/system/info');
      systemInfo = response.data;
      cpuPercentage = systemInfo.cpu.usage;
      memoryPercentage = Math.round((systemInfo.memory.used / systemInfo.memory.total) * 100);
      diskPercentage = Math.round((systemInfo.disk.used / systemInfo.disk.total) * 100);
    } catch (error) {
      console.error('Failed to fetch system info:', error);
    }
  }

  onMount(() => {
    fetchSystemInfo();
    const interval = setInterval(fetchSystemInfo, 5000);
    return () => clearInterval(interval);
  });
</script>

<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-8">Linux Web Panel</h1>
  
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
    <!-- 시스템 정보 -->
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">시스템 정보</h2>
      <div class="space-y-3">
        <p class="text-gray-600">호스트명: <span class="font-medium">{systemInfo.hostname}</span></p>
        <p class="text-gray-600">가동시간: <span class="font-medium">{systemInfo.uptime}</span></p>
        <p class="text-gray-600">운영체제: <span class="font-medium">{systemInfo.os}</span></p>
        <p class="text-gray-600">커널: <span class="font-medium">{systemInfo.kernel}</span></p>
      </div>
    </div>

    <!-- CPU 사용량 -->
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">CPU 사용량</h2>
      <div class="gauge-container">
        <div class="gauge">
          <div class="gauge-body" style="transform: rotate({45 + (cpuPercentage * 1.8)}deg)">
            <div class="gauge-cover">
              <div class="gauge-value" style="transform: rotate(-{45 + (cpuPercentage * 1.8)}deg)">{cpuPercentage}%</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 메모리 사용량 -->
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">메모리 사용량</h2>
      <div class="gauge-container">
        <div class="gauge">
          <div class="gauge-body memory" style="transform: rotate({45 + (memoryPercentage * 1.8)}deg)">
            <div class="gauge-cover">
              <div class="gauge-value" style="transform: rotate(-{45 + (memoryPercentage * 1.8)}deg)">{memoryPercentage}%</div>
            </div>
          </div>
        </div>
      </div>
      <div class="text-sm text-gray-600 mt-4">
        사용 중: {Math.round(systemInfo.memory.used / 1024 / 1024)}MB / 전체: {Math.round(systemInfo.memory.total / 1024 / 1024)}MB
      </div>
    </div>

    <!-- 디스크 사용량 -->
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">디스크 사용량</h2>
      <div class="gauge-container">
        <div class="gauge">
          <div class="gauge-body disk" style="transform: rotate({45 + (diskPercentage * 1.8)}deg)">
            <div class="gauge-cover">
              <div class="gauge-value" style="transform: rotate(-{45 + (diskPercentage * 1.8)}deg)">{diskPercentage}%</div>
            </div>
          </div>
        </div>
      </div>
      <div class="text-sm text-gray-600 mt-4">
        사용 중: {Math.round(systemInfo.disk.used / 1024 / 1024)}GB / 전체: {Math.round(systemInfo.disk.total / 1024 / 1024)}GB
      </div>
    </div>
  </div>

  <!-- 빠른 링크 -->
  <div class="mt-12">
    <h2 class="text-2xl font-bold mb-6">빠른 링크</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <a href="/overview" class="card hover:shadow-lg transition-shadow duration-200">
        <h3 class="text-lg font-semibold mb-2">시스템 개요</h3>
        <p class="text-gray-600">시스템 리소스 모니터링</p>
      </a>
      <a href="/shell" class="card hover:shadow-lg transition-shadow duration-200">
        <h3 class="text-lg font-semibold mb-2">웹 터미널</h3>
        <p class="text-gray-600">브라우저에서 직접 명령어 실행</p>
      </a>
      <a href="/daemons" class="card hover:shadow-lg transition-shadow duration-200">
        <h3 class="text-lg font-semibold mb-2">데몬 관리</h3>
        <p class="text-gray-600">시스템 서비스 관리</p>
      </a>
      <a href="/files" class="card hover:shadow-lg transition-shadow duration-200">
        <h3 class="text-lg font-semibold mb-2">파일 관리</h3>
        <p class="text-gray-600">파일 시스템 탐색 및 관리</p>
      </a>
    </div>
  </div>
</div>

<style>
  :global(body) {
    background-color: #f3f4f6;
  }

  .gauge-container {
    width: 200px;
    height: 200px;
    margin: 0 auto;
    position: relative;
  }

  .gauge {
    width: 100%;
    height: 100%;
    position: relative;
  }

  .gauge-body {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    border: 10px solid #f3f4f6;
    border-top-color: #3B82F6;
    transition: transform 1s ease-out;
    position: relative;
  }

  .gauge-body.memory {
    border-top-color: #10B981;
  }

  .gauge-body.disk {
    border-top-color: #8B5CF6;
  }

  .gauge-cover {
    width: 75%;
    height: 75%;
    background: white;
    border-radius: 50%;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
  }

  .gauge-value {
    font-size: 1.5rem;
    font-weight: bold;
    color: #1F2937;
  }
</style> 