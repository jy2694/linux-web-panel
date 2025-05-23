<script lang="ts">
  import { onMount } from 'svelte';
  import axios from 'axios';

  let systemInfo = {
    hostname: '',
    uptime: '',
    os: '',
    kernel: '',
    cpu: {
      usage: 0,
      cores: 0
    },
    memory: {
      total: 0,
      used: 0,
      free: 0
    },
    disk: {
      total: 0,
      used: 0,
      free: 0
    }
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
      console.error('시스템 정보를 가져오는데 실패했습니다:', error);
    }
  }

  onMount(() => {
    fetchSystemInfo();
    const interval = setInterval(fetchSystemInfo, 5000);
    return () => clearInterval(interval);
  });
</script>

<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-8">시스템 개요</h1>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
    <div class="card">
      <h2 class="text-xl font-semibold mb-6">시스템 정보</h2>
      <div class="space-y-4">
        <p class="flex justify-between">
          <span class="text-gray-600">호스트명:</span>
          <span class="font-medium">{systemInfo.hostname}</span>
        </p>
        <p class="flex justify-between">
          <span class="text-gray-600">가동 시간:</span>
          <span class="font-medium">{systemInfo.uptime}</span>
        </p>
        <p class="flex justify-between">
          <span class="text-gray-600">운영체제:</span>
          <span class="font-medium">{systemInfo.os}</span>
        </p>
        <p class="flex justify-between">
          <span class="text-gray-600">커널:</span>
          <span class="font-medium">{systemInfo.kernel}</span>
        </p>
      </div>
    </div>

    <div class="card">
      <h2 class="text-xl font-semibold mb-6">CPU 사용량</h2>
      <div class="space-y-4">
        <div class="flex justify-between">
          <span class="text-gray-600">사용률:</span>
          <span class="font-medium">{systemInfo.cpu.usage}%</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">코어 수:</span>
          <span class="font-medium">{systemInfo.cpu.cores}</span>
        </div>
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
    </div>

    <div class="card">
      <h2 class="text-xl font-semibold mb-6">메모리 사용량</h2>
      <div class="space-y-4">
        <div class="flex justify-between">
          <span class="text-gray-600">전체:</span>
          <span class="font-medium">{Math.round(systemInfo.memory.total / 1024 / 1024)} MB</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">사용 중:</span>
          <span class="font-medium">{Math.round(systemInfo.memory.used / 1024 / 1024)} MB</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">여유:</span>
          <span class="font-medium">{Math.round(systemInfo.memory.free / 1024 / 1024)} MB</span>
        </div>
        <div class="gauge-container">
          <div class="gauge">
            <div class="gauge-body memory" style="transform: rotate({45 + (memoryPercentage * 1.8)}deg)">
              <div class="gauge-cover">
                <div class="gauge-value" style="transform: rotate(-{45 + (memoryPercentage * 1.8)}deg)">{memoryPercentage}%</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <h2 class="text-xl font-semibold mb-6">디스크 사용량</h2>
      <div class="space-y-4">
        <div class="flex justify-between">
          <span class="text-gray-600">전체:</span>
          <span class="font-medium">{Math.round(systemInfo.disk.total / 1024 / 1024)} GB</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">사용 중:</span>
          <span class="font-medium">{Math.round(systemInfo.disk.used / 1024 / 1024)} GB</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">여유:</span>
          <span class="font-medium">{Math.round(systemInfo.disk.free / 1024 / 1024)} GB</span>
        </div>
        <div class="gauge-container">
          <div class="gauge">
            <div class="gauge-body disk" style="transform: rotate({45 + (diskPercentage * 1.8)}deg)">
              <div class="gauge-cover">
                <div class="gauge-value" style="transform: rotate(-{45 + (diskPercentage * 1.8)}deg)">{diskPercentage}%</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
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