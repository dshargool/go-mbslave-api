<script>
	/** @type {import('./$types').PageData}*/
	export let data;

	import { invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';

    function sortRegisters(multiplier = 1) {
        data.data.sort((a, b) => {
            const addressA = a.address.toString();
            const addressB = b.address.toString();
            
            // Split addresses on underscore
            const partsA = addressA.split('_');
            const partsB = addressB.split('_');
            
            // Get the first part (before underscore) as numbers
            const primaryA = parseInt(partsA[0]);
            const primaryB = parseInt(partsB[0]);
            
            // Primary sort: compare the first part numerically
            if (primaryA !== primaryB) {
                return (primaryA - primaryB) * multiplier;
            }
            
            // Secondary sort: if first parts are equal, compare second parts
            // If no second part exists, treat as 0 for sorting purposes
            const secondaryA = partsA.length > 1 ? parseInt(partsA[1]) || 0 : 0;
            const secondaryB = partsB.length > 1 ? parseInt(partsB[1]) || 0 : 0;
            
            return (secondaryA - secondaryB) * multiplier;
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
