<script lang="ts">
	import { Plus, Minus } from 'lucide-svelte';

	export let name: string;
	export let score: string;
	export let controls: boolean = false;
	export let orientation: 'left' | 'right' = 'left';
</script>

<div class="card variant-ghost flex font-mono text-6xl font-bold">
	{#if orientation === 'left'}
		<div class="grid w-full grid-cols-[auto_5rem]">
			<div class="flex w-full items-center justify-evenly">
				<div class="font-mono">{score}</div>
				<div>{name}</div>
			</div>
			<div class="btn-group variant-filled flex flex-col">
				<button class="variant-filled m-auto h-full w-full">
					<div class="flex justify-center">
						<Plus />
					</div>
				</button>
				<button class="variant-filled m-auto h-full w-full">
					<div class="flex justify-center">
						<Minus />
					</div>
				</button>
			</div>
		</div>
	{:else}
		<div class="grid w-full grid-cols-[5rem_auto]">
			<div class="btn-group variant-filled flex-col">
				<button class="m-auto h-full w-full">
					<div class="flex justify-center">
						<Plus />
					</div>
				</button>
				<button class="m-auto h-full w-full">
					<div class="flex justify-center">
						<Minus />
					</div>
				</button>
			</div>
			<div class="flex w-full items-center justify-evenly">
				<div>{name}</div>
				<div class="font-mono">{score}</div>
			</div>
		</div>
	{/if}
</div>
