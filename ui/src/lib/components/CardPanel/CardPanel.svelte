<script lang="ts">
	import { LiveScoreState } from '$lib/stores/HandballState';
	import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';

	$: lineups = $LiveScoreState && [
		{ players: $LiveScoreState.home.api.lineup, coach: $LiveScoreState.home.api.coach },
		{ players: $LiveScoreState.away.api.lineup, coach: $LiveScoreState.away.api.coach }
	];

	let selectedPerson: string = '';

	let actions = [
		{
			id: 'red',
			color: 'red',
			selectedColor: 'white',
			background: 'mistyrose',
			selectedBackground: 'red',
			border: 'red',
			text: 'Red Card'
		},
		{
			id: 'blue',
			color: 'blue',
			selectedColor: 'white',
			background: 'lightblue',
			selectedBackground: 'blue',
			border: 'blue',
			text: 'Blue Card'
		},
		{
			id: 'yellow',
			color: 'darkgoldenrod',
			selectedColor: 'black',
			background: 'lemonchiffon',
			selectedBackground: 'yellow',
			border: 'yellow',
			text: 'Yellow Card'
		}
	];
	let selectedAction: string = '';
	let playing: boolean = false;

	function play() {
		playing = true;
	}

	function stop() {
		playing = false;
	}
</script>

<div class="card variant-ghost grid w-full grid-rows-[6fr_1fr] gap-2">
	<div class="grid grid-cols-2 gap-2 text-white">
		{#if lineups}
			{#each lineups as l, i}
				<div class="flex p-2">
					<div
						class="scrollbar-none relative flex w-full flex-col justify-between overflow-y-auto pb-2"
					>
						{#each l.players as player, j}
							<div
								class="absolute w-full"
								style="top: {j * 2.5 + (j >= l.players.length ? 0.5 : 0)}rem;"
							>
								<!-- svelte-ignore a11y-click-events-have-key-events -->
								<label
									class="variant-ghost {selectedPerson === `player-${i}-${j}`
										? 'bg-red-800'
										: 'bg-red-600'} card flex cursor-pointer gap-2 hover:bg-red-900"
									for="player-{i}-{j}"
									on:click={() => (selectedPerson = `player-${i}-${j}`)}
								>
									<div class=" flex aspect-square h-9 rounded-full bg-red-700">
										<div class="m-auto">
											{player.number < 10 ? `0${player.number}` : `${player.number}`}
										</div>
									</div>
									<div class="flex flex-col justify-center">{player.firstname}</div>
									<div class="flex flex-col justify-center font-bold">{player.lastname}</div>
								</label>
								<input
									class="hidden"
									type="radio"
									bind:group={selectedPerson}
									value="player-{i}-{j}"
									id="player-{i}-{j}"
								/>
							</div>
						{/each}
						<div
							class="absolute w-full border-t-2 border-black pt-[0.625rem]"
							style="top: {l.players.length * 2.5 + 0.5}rem;"
						>
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<label
								class="variant-ghost {selectedPerson === `coach-${i}`
									? 'bg-red-800'
									: 'bg-red-600'} card flex cursor-pointer gap-2 hover:bg-red-900"
								for="player-{i}"
								on:click={() => (selectedPerson = `coach-${i}`)}
							>
								<div class=" flex aspect-square h-9 rounded-full bg-red-700">
									<div class="m-auto">
										{l.coach.letter}
									</div>
								</div>
								<div class="flex flex-col justify-center">{l.coach.firstname}</div>
								<div class="flex flex-col justify-center font-bold">{l.coach.lastname}</div>
							</label>
							<input
								class="hidden"
								type="radio"
								bind:group={selectedPerson}
								value="player-{i}"
								id="player-{i}"
							/>
						</div>
					</div>
				</div>
			{/each}
		{:else}
			<div class="flex w-full bg-yellow-500">
				<div
					class="scrollbar-none relative flex w-full flex-col justify-between overflow-y-auto pb-2"
				/>
			</div>
			<div class="flex w-full bg-yellow-500">
				<div
					class="scrollbar-none relative flex w-full flex-col justify-between overflow-y-auto pb-2"
				/>
			</div>
		{/if}
	</div>
	<div class="grid grid-cols-5 gap-2 p-2 text-lg">
		{#each actions as a}
			<div>
				<label
					class="flex h-full cursor-pointer border-2"
					for="action-{a.id}"
					style="
                        border-radius: var(--theme-rounded-base);
                        color: {selectedAction === a.id ? a.selectedColor : a.color};
                        font-weight: {selectedAction === a.id ? 'bold' : 'normal'};
                        border-color: {a.id};
                        background-color: {selectedAction === a.id
						? a.selectedBackground
						: a.background};
                        "
				>
					<div class="m-auto">{a.text}</div>
				</label>
				<input
					class="hidden"
					type="radio"
					bind:group={selectedAction}
					value={a.id}
					id="action-{a.id}"
				/>
			</div>
		{/each}
		<!-- <button class="bg-red-500 pl-2 hover:bg-red-600 active:bg-red-700">GO</button>
		<button class="bg-red-500 hover:bg-red-600 active:bg-red-700">STOP</button> -->
		<button
			on:click={play}
			disabled={playing}
			class="variant-filled-error btn border-2 border-black font-bold">GO</button
		>
		<button
			on:click={stop}
			disabled={!playing}
			class="variant-filled-error btn border-2 border-black font-bold">STOP</button
		>
	</div>
</div>

<style>
	.scrollbar-none::-webkit-scrollbar {
		display: none;
	}
</style>
