<script lang="ts">
	import type { AnyNode, Node } from '$lib/types';
	import { NodeType } from '$lib/types';
	import type { Component } from 'svelte';
	import EditLocationDialog from "./EditLocationDialog.svelte";
	import EditCodeDialog from "./EditCodeDialog.svelte";
	import EditRewardDialog from "./EditRewardDialog.svelte";
	import EditTimeDialog from "./EditTimeDialog.svelte";
	import DefaultEditDialog from "./DefaultEditDialog.svelte";

	type ComponentMap = { [K in NodeType]?: NodeFormComponent<Extract<AnyNode, { type: K }>> | NodeFormComponent<Node>; };
	type NodeFormComponent<N extends AnyNode> = Component<{
		node: N;
		onSave: (node: N) => void;
		onCancel: () => void;
	}>;

	let {
		node,
		onSave,
		onClose
	}: {
		node: AnyNode | null;
		onSave?: (node: AnyNode) => void;
		onClose?: () => void;
	} = $props();

	let dialog: HTMLDialogElement;

	const componentMap: ComponentMap = {
		[NodeType.LOCATION]: EditLocationDialog,
		[NodeType.CODE]: EditCodeDialog,
		[NodeType.MANUAL]: DefaultEditDialog,
		[NodeType.TIME]: EditTimeDialog,
		[NodeType.REWARD]: EditRewardDialog,
	};
	const BodyComponent = $derived(node ? componentMap[node.type] : undefined) as NodeFormComponent<AnyNode>;

	// Make dialog larger for reward nodes
	const isRewardNode = $derived(node?.type === NodeType.REWARD);

	$effect(() => {
		if (node) {
			dialog?.showModal();
		} else {
			dialog?.close();
		}
	});

	function handleSave(updatedNode: AnyNode) {
		onSave?.(updatedNode);
		onClose?.();
	}

	function handleCancel() {
		onClose?.();
	}
</script>

{#if node}
	<dialog bind:this={dialog} class="edit-dialog" class:large={isRewardNode} onclose={handleCancel} closedby="any">
		<div class="dialog-content">
			{#if BodyComponent}
				<BodyComponent node={node} onSave={handleSave} onCancel={handleCancel}></BodyComponent>
			{:else}
				<h2>⚠️ Cannot Edit Node</h2>
				<p>
					The node type <code>{node.type}</code> does not have an edit dialog implemented yet.
				</p>
				<div class="dialog-actions">
					<button type="button" class="btn-close" onclick={handleCancel}>
						Close
					</button>
				</div>
			{/if}
		</div>
	</dialog>
{/if}

<style>
	.edit-dialog {
		border: none;
		border-radius: 12px;
		padding: 0;
		max-width: 500px;
		width: 90%;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}

	.edit-dialog.large {
		max-width: 900px;
	}

	.edit-dialog::backdrop {
		background: rgba(0, 0, 0, 0.5);
	}

	.dialog-content {
		padding: 1.5rem;
	}

	.dialog-content h2 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.dialog-content p {
		margin: 0 0 1.5rem 0;
		color: #374151;
		line-height: 1.5;
	}

	code {
		background: #f3f4f6;
		padding: 0.125rem 0.375rem;
		border-radius: 0.25rem;
		font-family: monospace;
		font-size: 0.875em;
		color: #dc2626;
	}

	.dialog-actions {
		display: flex;
		justify-content: flex-end;
	}

	.btn-close {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		background: #dc2626;
		border: none;
		color: white;
	}

	.btn-close:hover {
		background: #b91c1c;
	}
</style>
