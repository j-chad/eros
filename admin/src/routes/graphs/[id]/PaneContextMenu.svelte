<script lang="ts">
	import {Gift, GitBranch, HandMetal, Key, MapPin, Play} from 'lucide-svelte';
    import {type AnyNode, type Node, NodeType} from "$lib/types";
	import {useNodes, useSvelteFlow} from "@xyflow/svelte";

	const { screenToFlowPosition } = useSvelteFlow();

    interface NodeOption {
        type: NodeType;
        label: string;
        icon: typeof Play;
        color: string;
        description: string;
    }

    let {
        x,
        y,
        onClose,
        onCreateNode
    }: {
        x: number;
        y: number;
        onClose?: () => void;
        onCreateNode?: (node: AnyNode) => void;
    } = $props();

    const nodeTypes: NodeOption[] = [
        {
            type: NodeType.LOCATION,
            label: 'Location',
            icon: MapPin,
            color: '#ec4899',
            description: 'GPS-based unlock'
        },
        {
            type: NodeType.CODE,
            label: 'Code',
            icon: Key,
            color: '#8b5cf6',
            description: 'Secret code unlock'
        },
        {
            type: NodeType.MANUAL,
            label: 'Manual',
            icon: HandMetal,
            color: '#10b981',
            description: 'Human confirmation'
        },
        {
            type: NodeType.REWARD,
            label: 'Reward',
            icon: Gift,
            color: '#06b6d4',
            description: 'Content reveal'
        }
    ];

    function handleCreateNode(type: NodeType) {
        onCreateNode?.({
            id: crypto.randomUUID(),
            type,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            title: "My New Node",
            ui_position: screenToFlowPosition({ x, y }),
        });
        onClose?.();
    }

    function handleClickOutside(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest('.context-menu')) {
            onClose?.();
        }
    }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="context-menu" style="left: {x}px; top: {y}px;">
    <div class="menu-header">
        <span class="menu-title">Create Node</span>
    </div>

    <div class="menu-items">
        {#each nodeTypes as nodeType}
            <button
                    class="menu-item"
                    onclick={() => handleCreateNode(nodeType.type)}
                    style="--node-color: {nodeType.color}"
            >
                <div class="item-icon" style="background: {nodeType.color}">
                    <nodeType.icon size={16} />
                </div>
                <div class="item-content">
                    <div class="item-label">{nodeType.label}</div>
                    <div class="item-description">{nodeType.description}</div>
                </div>
            </button>
        {/each}
    </div>
</div>

<style>
    .context-menu {
        position: fixed;
        background: white;
        border-radius: 8px;
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15), 0 0 0 1px rgba(0, 0, 0, 0.05);
        min-width: 240px;
        z-index: 1000;
        overflow: hidden;
        animation: fadeIn 0.15s ease-out;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: scale(0.95);
        }
        to {
            opacity: 1;
            transform: scale(1);
        }
    }

    .menu-header {
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #e5e7eb;
        background: #f9fafb;
    }

    .menu-title {
        font-size: 0.75rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: #6b7280;
    }

    .menu-items {
        padding: 0.5rem;
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .menu-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem;
        border: none;
        background: white;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.15s ease;
        text-align: left;
        width: 100%;
    }

    .menu-item:hover {
        background: #f3f4f6;
        transform: translateX(2px);
    }

    .menu-item:active {
        transform: translateX(2px) scale(0.98);
    }

    .item-icon {
        width: 36px;
        height: 36px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        flex-shrink: 0;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .item-content {
        flex: 1;
        min-width: 0;
    }

    .item-label {
        font-size: 0.875rem;
        font-weight: 600;
        color: #1f2937;
        margin-bottom: 0.125rem;
    }

    .item-description {
        font-size: 0.75rem;
        color: #6b7280;
        line-height: 1.3;
    }
</style>
