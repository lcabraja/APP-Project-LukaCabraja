<script>
	import { LiveScoreState } from '$lib/stores/HandballState';
	import { onMount } from 'svelte';

	let collapsed = new Set();

	const toggleCollapse = (path) => {
		if (collapsed.has(path)) {
			collapsed.delete(path);
		} else {
			collapsed.add(path);
		}
	};
</script>

<div class="json-browser">
	{#if $LiveScoreState}
		{#each Object.entries($LiveScoreState) as [key, value] (key)}
			<div class="json-item">
				<span class="json-key">{key}: </span>
				{#if typeof value === 'object' && value !== null}
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<span on:click={() => toggleCollapse(key)}>
						{#if collapsed.has(key)}
							{`{...}`}
						{:else}
							{`{`}
							<div class="json-object">
								{#each Object.entries(value) as [subKey, subValue] (`${key}.${subKey}`)}
									<div class="json-item">
										<span class="json-key">{subKey}: </span>
										{#if typeof subValue === 'object' && subValue !== null}
											<!-- svelte-ignore a11y-click-events-have-key-events -->
											<span on:click={() => toggleCollapse(`${key}.${subKey}`)}>
												{#if collapsed.has(`${key}.${subKey}`)}
													{`{...}`}
												{:else}
													{`{`}
													<div class="json-object">
														{#each Object.entries(subValue) as [innerSubKey, innerSubValue] (`${key}.${subKey}.${innerSubKey}`)}
															<div class="json-item">
																<span class="json-key">{innerSubKey}: </span>
																{#if typeof innerSubValue === 'object' && innerSubValue !== null}
																	<span
																		on:click={() =>
																			toggleCollapse(`${key}.${subKey}.${innerSubKey}`)}
																	>
																		{#if collapsed.has(`${key}.${subKey}.${innerSubKey}`)}
																			{`{...}`}
																		{:else}
																			{JSON.stringify(innerSubValue, null, 2)}
																		{/if}
																	</span>
																{:else}
																	{JSON.stringify(innerSubValue)}
																{/if}
															</div>
														{/each}
													</div>
													{`}`}
												{/if}
											</span>
										{:else}
											{JSON.stringify(subValue)}
										{/if}
									</div>
								{/each}
							</div>
							{`}`}
						{/if}
					</span>
				{:else}
					{JSON.stringify(value)}
				{/if}
			</div>
		{/each}
	{/if}
</div>

<style>
	.json-key {
		font-weight: bold;
	}
	.json-object {
		margin-left: 20px;
	}
</style>
