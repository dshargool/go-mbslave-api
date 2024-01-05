/** @type {import('./$types').PageLoad} */
export async function load({ fetch, params }) {
	const resp = await fetch('http://127.0.0.1:8081/all_registers');
	if (resp.ok) {
		const data = await resp.json();
		return { data: data };
	} else {
		console.error('Failed to fetch data from PI');
		return {
			status: resp.status,
			error: new Error('Failed to fetch data')
		};
	}
}
