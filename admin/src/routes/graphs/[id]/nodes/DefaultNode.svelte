<script lang="ts">
	import {type NodeProps} from '@xyflow/svelte';
	import {MessageSquareWarning} from 'lucide-svelte';

	let {}: NodeProps = $props();
</script>

<div class="default-node">
	<div class="icon-wrapper">
		<MessageSquareWarning size={16} />
	</div>
	<div class="title">
		<div class="name">Unknown Node</div>
	</div>
</div>

<style>
	.default-node {
		background: #ef4444;
		border: 2px solid #b91c1c;
		border-radius: 8px;
		min-width: 200px;
		max-width: 280px;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);

		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		background: linear-gradient(135deg, #ef4444 0%, #b91c1c 100%);
		color: white;
	}

	.icon-wrapper {
		width: 32px;
		height: 32px;
		border-radius: 6px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.2);
	}

	.title {
		flex: 1;
		min-width: 0;
	}

	.name {
		font-size: 0.875rem;
		font-weight: 600;
		margin-top: 0.125rem;
	}
</style>
