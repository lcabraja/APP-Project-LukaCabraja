<script lang="ts">
	import Penalties from '$lib/components/Blocks/Penalties.svelte';
	import Score from '$lib/components/Blocks/Score.svelte';
	import StatisticsPanel from '$lib/components/Blocks/StatisticsPanel.svelte';
	import Timer from '$lib/components/Blocks/Timer.svelte';
	import CardPanel from '$lib/components/CardPanel/CardPanel.svelte';
	import type { StopWatch } from '$lib/data/StopWatch';

	let homeScore = '00';
	let awayScore = '00';
	let teamHome = 'HOM';
	let teamAway = 'GUE';
	let colorHome = 'blue';
	let colorAway = 'red';

	let mainStopWatch: StopWatch;
	let homeTimeouts: Array<string | undefined> = ['2:00.000', undefined];
	let awayTimeouts: Array<string> = ['1:23.456', '1:32.456'];

	export let defaultOrientation: 'left' | 'right' = 'right';
	function inverseOrientation(): 'left' | 'right' {
		return defaultOrientation === 'left' ? 'right' : 'left';
	}

	let lineup = Array.from({ length: 2 }).map((_, i) =>
		Array.from({ length: 16 }).map((_, j) => ({
			number: j + 1,
			firstName: 'Ime-' + (j + 1),
			lastName: 'Prezime-' + (j + 1)
		}))
	);
	let selectedPlayer = '';

	let selectedAction = '';

	let playing: NodeJS.Timeout;
	const play = () => {
		playing = setTimeout(() => {
			selectedPlayer = '';
			selectedAction = '';
		}, 2500);
	};
	const stop = () => {
		clearTimeout(playing);
		selectedPlayer = '';
		selectedAction = '';
	};
</script>

<div class="flex h-full w-full gap-2 p-8">
	<div class="grid w-full grid-rows-[1fr_1fr_3fr] gap-2">
		<div class="grid grid-cols-2 gap-2">
			<Score name={teamHome} score={homeScore} orientation={defaultOrientation} controls={true} />
			<Score name={teamAway} score={awayScore} orientation={inverseOrientation()} controls={true} />
		</div>
		<div class="card variant-ghost flex justify-between gap-2 bg-green-500 p-2">
			<Penalties
				times={['0:00', '0:00', '0:00']}
				orientation={defaultOrientation}
				teamName={teamHome}
			/>
			<Timer time="00:00:000" controls={true} />
			<Penalties
				times={['0:00', '0:00', '0:00']}
				orientation={inverseOrientation()}
				teamName={teamAway}
			/>
		</div>
		<!-- <div class="grid grid-cols-2 grid-rows-4 gap-2">
			<div class="card variant-ghost">A</div>
			<div class="card variant-ghost">B</div>
			<div class="card variant-ghost">C</div>
			<div class="card variant-ghost">D</div>
			<div class="card variant-ghost">E</div>
			<div class="card variant-ghost">F</div>
			<div class="card variant-ghost">G</div>
			<div class="card variant-ghost">H</div>
		</div> -->
		<div class="grid grid-cols-2 gap-2">
			<StatisticsPanel orientation={defaultOrientation} teamName={teamHome} />
			<StatisticsPanel orientation={inverseOrientation()} teamName={teamAway} />
		</div>
	</div>
	<CardPanel />
</div>
