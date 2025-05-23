declare module 'svelte-apexcharts' {
  import { SvelteComponent } from 'svelte';

  interface ApexChartsProps {
    options: any;
  }

  export default class ApexCharts extends SvelteComponent<ApexChartsProps> {}
} 