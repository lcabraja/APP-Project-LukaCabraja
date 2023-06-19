<script lang="ts">
	import '@skeletonlabs/skeleton/themes/theme-seafoam.css';
	import '@skeletonlabs/skeleton/styles/all.css';
	import '../app.postcss';

	import { page } from '$app/stores';
	import {
		AppBar,
		AppRail,
		AppRailAnchor,
		AppShell,
		autoModeWatcher,
		Avatar
	} from '@skeletonlabs/skeleton';
	import { LucideScreenShare, Save, Settings, Timer } from 'lucide-svelte';
	import { onMount } from 'svelte';

	import { createWorker } from '$lib/api/worker';
	import { Action } from '$lib/api/ws';
	import {
		LiveScoreState,
		PreviousLiveScoreState,
		subscribeHandler
	} from '$lib/stores/HandballState';
	import { Connection } from '$lib/stores/ConnectionState';

	let tbMatchid: HTMLInputElement;

	$: if ($LiveScoreState) {
		tbMatchid.value = $LiveScoreState.matchid;
	}

	onMount(() => {
		createWorker().port.postMessage(['send', JSON.stringify(Action.onboard)]);
		LiveScoreState.subscribe((cs) => subscribeHandler(cs, $PreviousLiveScoreState));
	});
</script>

<svelte:head>
	{@html `<script>${autoModeWatcher.toString()} autoModeWatcher();</script>`}
</svelte:head>

<AppShell>
	<svelte:fragment slot="header">
		<AppBar>
			<svelte:fragment slot="lead">
				<!-- Status Icon -->
				<Avatar
					initials={$Connection.initials}
					fill={$Connection.fill}
					width="w-12"
					background={$Connection.background}
				/>
			</svelte:fragment>
			<svelte:fragment slot="trail">
				<div class="input-group-divider input-group grid-cols-[auto_1fr_auto]">
					<div class="input-group-shim">Match ID</div>
					<input class="w-[60px]" bind:this={tbMatchid} type="text" />
					<button
						on:click={() => {
							if ($LiveScoreState) $LiveScoreState.matchid = tbMatchid.value;
						}}
						class="variant-filled-secondary active:variant-filled-primary"
					>
						<Save />
					</button>
				</div>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>
	<svelte:fragment slot="sidebarLeft">
		<AppRail>
			<AppRailAnchor href="/" selected={$page.url.pathname === '/'} title="Score & Timer">
				<svelte:fragment slot="lead"><Timer class="m-auto" /></svelte:fragment>
				<span>SCORE</span>
			</AppRailAnchor>
			<AppRailAnchor href="/gfx" selected={$page.url.pathname === '/gfx'} title="Graphics">
				<svelte:fragment slot="lead"><LucideScreenShare class="m-auto" /></svelte:fragment>
				<span>GFX</span>
			</AppRailAnchor>
			<AppRailAnchor
				href="/switchboard"
				selected={$page.url.pathname === '/switchboard'}
				title="Switchboard"
			>
				<svelte:fragment slot="lead"><Settings class="m-auto" /></svelte:fragment>
				<span>SWB</span>
			</AppRailAnchor>
		</AppRail>
	</svelte:fragment>
	<slot />
</AppShell>
