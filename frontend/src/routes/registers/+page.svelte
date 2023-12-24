<script>
	/** @type {import('./$types').PageData}*/
	export let data;

    function sortRegisters(multiplier) {
        data.data.sort((a,b) => { 
            if (a.address > b.address) {
                return 1 * multiplier;
            }
            if (a.address < b.address) {
                return -1 * multiplier;
            } 
            return 0;
        })
        data.data = data.data; // Refresh data
    }
	import { invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';

	function rerunLoadFunction() {
		// any of these will cause the `load` function to rerun
		invalidateAll();
	}

    onMount(async () => {
        setInterval(() => {
            rerunLoadFunction();
            sortRegisters(1);
        }, 10000);
    });



</script>

<main>
	<h1 class="font-bold text-2xl">Registers</h1>
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
