<script>
	/** @type {import('./$types').PageData}*/
	export let data;

	import { invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';

	function sortRegisters(multiplier = 1) {
		data.data.sort((a, b) => {
			if (a.address > b.address) {
				return 1 * multiplier;
			}
			if (a.address < b.address) {
				return -1 * multiplier;
			}
			return 0;
		});
		data.data = data.data; // Refresh data
	}

	onMount(async () => {
	    sortRegisters(1);
		setInterval(() => {
			invalidateAll().then(() => {
				sortRegisters(1);
			});
		}, 10000);
	});

	// Start doing things
</script>

<main>
	<h3>Registers</h3>
	<div class="overflow-x-auto">
		<table class="table table-zebra table-pin-rows">
			<thead>
				<tr>
					<th>Tag</th>
					<th>Description</th>
					<th>Address</th>
					<th>Value</th>
					<th>Last Update</th>
				</tr>
			</thead>
			<tbody>
				{#if data.data}
					{#each data.data as register}
						<tr>
							<td>{register.tag}</td>
							<td>{register.description}</td>
							<td>{register.address}</td>
							<td>{register.value}</td>
							<td>{register.last_update}</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</main>
