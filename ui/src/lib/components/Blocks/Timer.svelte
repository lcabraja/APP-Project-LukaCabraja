<script lang="ts">
	import { Play, Pause } from 'lucide-svelte';
	import { each } from 'svelte/internal';

	export let time: string;
	export let timeRunning: boolean = false;
	export let controls: boolean = false;

	$: segments = time.split(':');
</script>

<div class="flex flex-col items-center justify-between text-6xl">
	<div class="flex h-full w-full items-center font-mono">
		<span class="m-auto">
			{#each segments as segment, i}
				{#if i < segments.length - 1}
					{segment}<span class="mx-[-0.6rem] {timeRunning && 'blink'}">:</span>
				{:else}
					{segment}
				{/if}
			{/each}
		</span>
	</div>
	{#if controls}
		<div class="flex w-full gap-2 align-bottom">
			<div class="btn-group variant-filled mx-auto h-12">
				<button disabled={timeRunning}><Play /></button>
				<button disabled={!timeRunning}><Pause /></button>
			</div>
			<div class="btn-group variant-filled mx-auto h-12">
				<button>+1s</button>
				<button>-1s</button>
			</div>
		</div>
	{/if}
</div>

<style>
	.blink {
		animation: blink 2s infinite;
	}

	@keyframes blink {
		0% {
			opacity: 0;
		}
		49% {
			opacity: 0;
		}
		50% {
			opacity: 1;
		}
		99% {
			opacity: 1;
		}
		100% {
			opacity: 0;
		}
	}
</style>
