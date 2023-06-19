<script lang="ts">
	import { Action } from '$lib/api/ws';
	import { LiveScoreState } from '$lib/stores/HandballState';
	import { SendAction } from '$lib/stores/WebSocket';

	export let orientation: 'left' | 'right' = 'left';
	export let teamName: string;
	$: team = teamName.slice(0, 3).toUpperCase();

	let elements = [
		{
			name: 'Scoring',
			local: orientation === 'right' ? $LiveScoreState?.home.local.scoring.percent : 'away',
			api: 'red',
			add: () => {
				$SendAction &&
					$SendAction(
						Action.offset_int[orientation === 'right' ? 'home' : 'away'].score_attempt_add1
					);
			},
			sub: () => {
				$SendAction &&
					$SendAction(
						Action.offset_int[orientation === 'right' ? 'home' : 'away'].score_attempt_sub1
					);
			}
		},
		{
			name: 'Defense',
			local: 'red',
			api: 'red',
			add: () => {
				$SendAction &&
					$SendAction(
						Action.offset_int[orientation === 'right' ? 'home' : 'away'].defense_success_add1
					);
			},
			sub: () => {
				$SendAction &&
					$SendAction(
						Action.offset_int[orientation === 'right' ? 'home' : 'away'].defense_success_sub1
					);
			}
		},
		{
			name: 'FBG',
			local: 'red',
			api: 'red',
			add: () => {
				$SendAction &&
					$SendAction(Action.offset_int[orientation === 'right' ? 'home' : 'away'].fbg_add1);
			},
			sub: () => {
				$SendAction &&
					$SendAction(Action.offset_int[orientation === 'right' ? 'home' : 'away'].fbg_sub1);
			}
		},
		{
			name: '7m Goals',
			local: 'red',
			api: 'red',
			add: () => {
				$SendAction &&
					$SendAction(Action.offset_int[orientation === 'right' ? 'home' : 'away'].goal7m_add1);
			},
			sub: () => {
				$SendAction &&
					$SendAction(Action.offset_int[orientation === 'right' ? 'home' : 'away'].goal7m_sub1);
			}
		}
	];
</script>

<div class="card variant-ghost">
	<div class="flex h-full flex-col justify-center">
		<div class="flex h-10 justify-center p-2">
			<span class="text-xl"><b>{team}</b>&nbsp;Statistics</span>
		</div>
		<hr />
		{#if elements}
			{#each elements as { name, local, api, add, sub }, i}
				<div class="flex h-full gap-2 p-2">
					<span class="vertical border-primary cursor-vertical-text border-r-[1px] text-center"
						>{name}</span
					>
					<div class="grid w-full grid-cols-[auto_1fr] grid-rows-2 items-center">
						<span class="pr-6">Local:</span>
						<span>{local}</span>
						<span class="pr-6">API:</span>
						<span>{api}</span>
					</div>
					<div class="btn-group variant-filled flex flex-col text-xl">
						<button on:click={add} class="h-full border-b-[1px] border-black">+</button>
						<button on:click={sub} class="h-full">-</button>
					</div>
				</div>
				{#if i < elements.length - 1}
					<hr />
				{/if}
			{/each}
		{/if}
		<!-- 
		<div class="flex h-full p-2">
			<span class="vertical text-center cursor-vertical-text">Scoring</span>
			<div class="bg-red-500">Local</div>
		</div>
		<hr />
		<div class="flex h-full p-2">Defense</div>
		<hr />
		<div class="flex h-full p-2">Fast Break Goals</div>
		<hr />
		<div class="flex h-full p-2">7m Goals</div> -->
	</div>
</div>

<style>
	.vertical {
		writing-mode: sideways-lr;
		text-orientation: sideways;
	}
</style>
