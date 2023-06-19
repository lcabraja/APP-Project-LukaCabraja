<script lang="ts">
	import { Action } from '$lib/api/ws';
	import { SendAction } from '$lib/stores/WebSocket';
</script>

<div class="flex h-full w-full">
	<button
		class="btn variant-filled-primary m-auto text-xl"
		on:click={() => $SendAction && $SendAction(Action.resume)}>Send</button
	>
</div>
