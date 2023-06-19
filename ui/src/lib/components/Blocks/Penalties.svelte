<script lang="ts">
	import { X } from 'lucide-svelte';

	export let times: string[] = [];
	export let orientation: 'left' | 'right' = 'left';
	export let teamName: string;

	$: team = teamName.slice(0, 3).toUpperCase();
</script>

<div class="flex flex-row gap-2">
	{#if orientation === 'left'}
		<div class="flex w-full">
			<button class="btn variant-filled">
				<div class="flex flex-col">
					<span>2 min</span>
					<b><span>{team}</span></b>
				</div>
			</button>
		</div>
		<div class="flex flex-col justify-between gap-2 text-2xl">
			{#each times as time}
				<div class="btn-group flex h-full flex-row">
					<div class="input-group-shim variant-filled flex items-center px-5">{time}</div>
					<button class="variant-filled text-primary-500"><X /></button>
				</div>
			{/each}
		</div>
	{:else}
		<div class="flex flex-col justify-between gap-2 text-2xl">
			{#each times as time}
				<div class="btn-group flex h-full flex-row">
					<button class="variant-filled text-primary-500"><X /></button>
					<div class="input-group-shim variant-filled flex items-center px-5">{time}</div>
				</div>
			{/each}
		</div>
		<div class="flex w-full">
			<button class="btn variant-filled">
				<div class="flex flex-col">
					<span>2 min</span>
					<b><span>{team}</span></b>
				</div>
			</button>
		</div>
	{/if}
</div>
