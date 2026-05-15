<script lang="ts">
	import { goto } from '$app/navigation';
	import { ChevronLeft, Heart, Check } from 'lucide-svelte';
	import BrandHeader from '$lib/ui/BrandHeader.svelte';
	import Card from '$lib/ui/base/Card.svelte';
	import StartNodeView from '$lib/ui/nodes/StartNode.svelte';
	import CodeGateView from '$lib/ui/nodes/CodeGate.svelte';
	import LocationGateView from '$lib/ui/nodes/LocationGate.svelte';
	import ManualGateView from '$lib/ui/nodes/ManualGate.svelte';
	import TimeGateView from '$lib/ui/nodes/TimeGate.svelte';
	import RewardNodeView from '$lib/ui/nodes/RewardNode.svelte';
	import type { AnyNode, GraphDetail } from '$lib/types/graph';
	import type { UnlockResult } from '$lib/api/graph.api';
	import { NodeType } from '$lib/types/graph';
	import Confetti from '$lib/ui/Confetti.svelte';

	const { data }: { data: { graph: GraphDetail } } = $props();

	// Make nodes and edges reactive and mutable
	let nodes = $state(data.graph.nodes ?? []);
	let edges = $state(data.graph.edges ?? []);

	// Celebration state
	let celebrationState = $state<'gate-success' | 'reward-reveal' | null>(null);
	let celebrationNode = $state<AnyNode | null>(null);

	// --- current node resolution ---
	//
	// The backend returns only accessible nodes (unlocked + one hop ahead).
	// We determine the "current" node as follows:
	//   1. Find all unlocked nodes (unlocked_at is set, or type is START).
	//   2. Among accessible (returned) nodes that are NOT unlocked, those are
	//      the frontier - nodes available to attempt next.
	//   3. If there is exactly one frontier node, that becomes the active node.
	//   4. If there are multiple frontier nodes (branching path), show a choice screen.
	//   5. If there are no frontier nodes, the graph is complete - show the last
	//      unlocked node (likely a reward).

	const unlockedNodes = $derived(
		nodes.filter((n) => n.unlocked_at != null || n.type === NodeType.START),
	);

	const unlockedIds = $derived(new Set(unlockedNodes.map((n) => n.id)));

	const frontierNodes = $derived(nodes.filter((n) => !unlockedIds.has(n.id)));

	// The "latest" unlocked node: most recently unlocked by timestamp, falling back to the start node.
	const latestUnlocked = $derived((): AnyNode | null => {
		if (unlockedNodes.length === 0) return null;
		return unlockedNodes.reduce((best, n) => {
			if (!best) return n;
			if (!best.unlocked_at) return n;
			if (!n.unlocked_at) return best;
			return n.unlocked_at >= best.unlocked_at ? n : best;
		});
	});

	// Outgoing edges from the latest unlocked node, to determine if there are branching choices.
	const outgoingEdges = $derived(
		latestUnlocked() ? edges.filter((e) => e.from === latestUnlocked()!.id) : [],
	);

	// Determine what to show.
	// - If there are frontier nodes reachable from the latest unlocked → show them (active gate/reward)
	// - If no frontier (complete) → show the last unlocked node
	const activeNodes = $derived((): AnyNode[] => {
		// Frontier nodes that are directly reachable from the latest unlocked node.
		const directFrontier = frontierNodes.filter((n) =>
			outgoingEdges.some((e) => e.to === n.id),
		);
		if (directFrontier.length > 0) return directFrontier;

		// No reachable frontier - graph complete.
		// Use graph topology: show terminal node(s) - nodes with no outgoing edges.
		// This avoids timestamp tie-breaking issues when the backend auto-unlocks
		// reward nodes at the same instant as their preceding gate.
		const nodesWithOutgoing = new Set(edges.map((e) => e.from));
		const terminalNodes = unlockedNodes.filter((n) => !nodesWithOutgoing.has(n.id));
		if (terminalNodes.length === 1) return terminalNodes;
		if (terminalNodes.length > 1) {
			// Multiple terminals (branching endings): pick the most recently unlocked.
			return [terminalNodes.reduce((best, n) => {
				if (!best.unlocked_at) return n;
				if (!n.unlocked_at) return best;
				return n.unlocked_at > best.unlocked_at ? n : best;
			})];
		}

		// Fallback (shouldn't normally be reached).
		const last = latestUnlocked();
		return last ? [last] : [];
	});

	// When there's more than one active node, the user must pick a path.
	const isBranching = $derived(activeNodes().length > 1);

	// The single node being displayed (null when branching or empty).
	const currentNode = $derived(activeNodes().length === 1 ? activeNodes()[0] : null);

	// Branch choice: which node did the user pick (only relevant when isBranching).
	let chosenNodeId = $state<string | null>(null);
	const chosenNode = $derived(
		chosenNodeId ? nodes.find((n) => n.id === chosenNodeId) ?? null : null,
	);

	// The node to actually render: explicit choice > single current.
	const displayNode = $derived(chosenNode ?? currentNode);

	// Edge label for the chosen branch (used in the back button label).
	const chosenEdgeLabel = $derived(
		chosenNodeId
			? (outgoingEdges.find((e) => e.to === chosenNodeId)?.choice_label ?? null)
			: null,
	);

	function handleChoose(nodeId: string) {
		chosenNodeId = nodeId;
	}

	function handleBack() {
		if (chosenNodeId) {
			chosenNodeId = null;
		} else {
			goto('/');
		}
	}

	function handleUnlock(result: UnlockResult) {
		// Reset branch choice so we advance past the completed gate.
		chosenNodeId = null;

		// Update the unlocked node
		const nodeIndex = nodes.findIndex(n => n.id === result.unlocked_node.id);
		if (nodeIndex >= 0) {
			nodes[nodeIndex] = result.unlocked_node;
		}

		// Add new nodes and edges
		nodes.push(...result.new_nodes);
		edges.push(...result.new_edges);

		// Determine celebration type
		const hasRewardInNew = result.new_nodes.some(n => n.type === NodeType.REWARD);
		celebrationState = hasRewardInNew ? 'reward-reveal' : 'gate-success';
		celebrationNode = result.unlocked_node;

		// Auto-advance after celebration
		const delay = hasRewardInNew ? 3500 : 1500;
		setTimeout(() => {
			celebrationState = null;
			celebrationNode = null;
		}, delay);
	}
</script>

<svelte:head>
	<title>{data.graph.title} - Eros</title>
</svelte:head>

<div class="mx-auto min-h-dvh max-w-md px-4 py-6 flex flex-col gap-6">
	<!-- Header -->
	<div class="flex items-center gap-3">
		<button
			class="btn btn-ghost btn-sm btn-circle"
			onclick={handleBack}
			aria-label={chosenNodeId ? 'Back to choices' : 'Back'}
		>
			<ChevronLeft size={20} />
		</button>
		<div class="flex-1 min-w-0">
			<h2 class="text-sm font-bold truncate">{data.graph.title}</h2>
			{#if chosenEdgeLabel}
				<p class="text-xs opacity-50 truncate">{chosenEdgeLabel}</p>
			{/if}
		</div>
		<BrandHeader compact />
	</div>

	<!-- Main content -->
	<div class="flex-1 flex flex-col justify-center" style="view-transition-name: card-content">
		{#if nodes.length === 0}
			<!-- Graph has no accessible nodes yet -->
			<Card>
				<div class="flex flex-col items-center gap-3 py-6 text-center animate-popIn">
					<p class="text-base font-semibold">Nothing to see yet.</p>
					<p class="text-sm opacity-60">Check back once the adventure has started.</p>
				</div>
			</Card>

		{:else if isBranching && !chosenNodeId}
			<!-- Branch choice screen -->
			<div class="flex flex-col gap-4 animate-popIn">
				<div class="text-center">
					<p class="text-xs font-semibold opacity-60 uppercase tracking-wide mb-2">Choose your path</p>
					<h1 class="text-xl font-extrabold">{latestUnlocked()?.title ?? 'Next step'}</h1>
				</div>
				<div class="flex flex-col gap-3">
					{#each activeNodes() as node (node.id)}
						{@const edgeLabel = outgoingEdges.find((e) => e.to === node.id)?.choice_label}
						<button
							onclick={() => handleChoose(node.id)}
							class="w-full text-left"
						>
							<Card>
								<div class="flex flex-col gap-1 py-1">
									<p class="font-semibold">{edgeLabel ?? node.title}</p>
									{#if node.description}
										<p class="text-sm opacity-60 line-clamp-2">{node.description}</p>
									{/if}
								</div>
							</Card>
						</button>
					{/each}
				</div>
			</div>

		{:else if displayNode}
			<!-- Single node view -->
			<Card>
				<div
					class="relative py-4 px-2"
					class:animate-celebrateShake={celebrationState === 'reward-reveal'}
				>
					{#if celebrationState === 'reward-reveal'}
						<Confetti duration={3500} />
					{/if}
					{#if celebrationState}
						<!-- Celebration overlay -->
						<div class="flex flex-col items-center gap-6 text-center"
							 class:animate-rewardReveal={celebrationState === 'reward-reveal'}
							 class:animate-popIn={celebrationState === 'gate-success'}
						>
							<div
								class="w-20 h-20 rounded-full flex items-center justify-center shadow-lg {celebrationState === 'reward-reveal' ? 'bg-primary/15 shadow-primary/30 animate-glowPulse' : 'bg-success/15'}"
							>
								{#if celebrationState === 'reward-reveal'}
									<Heart size={32} class="text-primary" />
								{:else}
									<Check size={32} class="text-success" />
								{/if}
							</div>
							<div class="flex flex-col gap-2">
								<p class="text-xs font-semibold opacity-60 uppercase tracking-wide">
									{celebrationState === 'reward-reveal' ? 'You unlocked something special' : 'Unlocked!'}
								</p>
								<h1 class="text-2xl font-extrabold">{celebrationNode?.title ?? 'Success'}</h1>
							</div>
						</div>
					{:else if displayNode.type === NodeType.START}
						<StartNodeView node={displayNode} />
					{:else if displayNode.type === NodeType.CODE}
						<CodeGateView node={displayNode} graphId={data.graph.id} onUnlock={handleUnlock} />
					{:else if displayNode.type === NodeType.LOCATION}
						<LocationGateView node={displayNode} graphId={data.graph.id} onUnlock={handleUnlock} />
					{:else if displayNode.type === NodeType.MANUAL}
						<ManualGateView node={displayNode} graphId={data.graph.id} onUnlock={handleUnlock} />
					{:else if displayNode.type === NodeType.TIME}
						<TimeGateView node={displayNode} graphId={data.graph.id} onUnlock={handleUnlock} />
					{:else if displayNode.type === NodeType.REWARD}
						<RewardNodeView node={displayNode} />
					{/if}
				</div>
			</Card>

		{:else}
			<!-- Empty/error state -->
			<Card>
				<div class="flex flex-col items-center gap-3 py-6 text-center animate-popIn">
					<p class="text-sm opacity-60">Something went wrong loading this adventure.</p>
				</div>
			</Card>
		{/if}
	</div>
</div>
