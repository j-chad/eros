<script lang="ts">
    import type {StartNode} from "$lib/types";
    import {Edit, Plus, Trash2} from "lucide-svelte";

    let { nodes, day, date, isCurrentMonth, onDeleteFlow, onCreateFlow, onEditFlow }: {
        nodes: StartNode[];
        day: number;
        date: Date;
        isCurrentMonth: boolean;
        onCreateFlow?: (date: Date) => Promise<void>;
        onEditFlow?: (id: string) => Promise<void>;
        onDeleteFlow?: (id: string) => Promise<void>;
    } = $props();

    function isToday(date: Date): boolean {
        const today = new Date();
        return date.toDateString() === today.toDateString();
    }
</script>

<div
        class="calendar-day"
        class:other-month={!isCurrentMonth}
        class:today={isToday(date)}
>
    <div class="day-number">{day}</div>

    {#if isCurrentMonth}
        {#if nodes.length > 0}
            <div class="flows">
                {#each nodes as flow}
                    <div class="flow-item">
                        <div class="flow-content">
                            <div class="flow-title">{flow.title}</div>
                            {#if flow.description}
                                <div class="flow-description">{flow.description}</div>
                            {/if}
                        </div>
                        <div class="flow-actions">
                            <button
                                    class="action-btn edit"
                                    onclick={() => onEditFlow?.(flow.id)}
                                    title="Edit flow"
                            >
                                <Edit size={12} />
                            </button>
                            <button
                                    class="action-btn delete"
                                    onclick={() => onDeleteFlow?.(flow.id)}
                                    title="Delete flow"
                            >
                                <Trash2 size={12} />
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}

        <button
                class="add-flow-btn"
                onclick={() => onCreateFlow?.(date)}
                title="Create flow"
        >
            <Plus size={14} />
        </button>
    {/if}
</div>

<style>
    .calendar-day {
        background: white;
        min-height: 120px;
        padding: 0.5rem;
        position: relative;
        display: flex;
        flex-direction: column;
    }

    .calendar-day.other-month {
        background: #f9fafb;
        opacity: 0.5;
    }

    .calendar-day.today {
        background: #eff6ff;
    }

    .day-number {
        font-size: 0.875rem;
        font-weight: 500;
        color: #374151;
        margin-bottom: 0.5rem;
    }

    .calendar-day.today .day-number {
        background: #3b82f6;
        color: white;
        width: 1.75rem;
        height: 1.75rem;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
    }

    .flows {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        flex: 1;
        overflow-y: auto;
    }

    .flow-item {
        background: #f0f9ff;
        border: 1px solid #bae6fd;
        border-radius: 4px;
        padding: 0.5rem;
        font-size: 0.75rem;
        display: flex;
        justify-content: space-between;
        gap: 0.5rem;
        transition: all 0.2s;
    }

    .flow-item:hover {
        background: #e0f2fe;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .flow-content {
        flex: 1;
        min-width: 0;
    }

    .flow-title {
        font-weight: 600;
        color: #0c4a6e;
        margin-bottom: 0.25rem;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .flow-description {
        color: #475569;
        font-size: 0.6875rem;
        line-height: 1.3;
        margin-bottom: 0.25rem;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .flow-time {
        color: #64748b;
        font-size: 0.6875rem;
    }

    .flow-actions {
        display: flex;
        gap: 0.25rem;
        flex-shrink: 0;
    }

    .action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 1.5rem;
        height: 1.5rem;
        border: none;
        border-radius: 3px;
        cursor: pointer;
        transition: all 0.2s;
        opacity: 0.7;
    }

    .action-btn:hover {
        opacity: 1;
    }

    .action-btn.edit {
        background: #dbeafe;
        color: #1e40af;
    }

    .action-btn.edit:hover {
        background: #bfdbfe;
    }

    .action-btn.delete {
        background: #fee2e2;
        color: #991b1b;
    }

    .action-btn.delete:hover {
        background: #fecaca;
    }

    .add-flow-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        padding: 0.5rem;
        margin-top: 0.5rem;
        background: #f3f4f6;
        border: 1px dashed #d1d5db;
        border-radius: 4px;
        color: #6b7280;
        cursor: pointer;
        transition: all 0.2s;
        font-size: 0.75rem;
    }

    .add-flow-btn:hover {
        background: #e5e7eb;
        border-color: #9ca3af;
        color: #374151;
    }

    @media (max-width: 1024px) {
        .calendar-day {
            min-height: 100px;
        }

        .flow-description {
            display: none;
        }
    }

    @media (max-width: 768px) {
        .calendar-day {
            min-height: 80px;
            padding: 0.25rem;
        }

        .day-number {
            font-size: 0.75rem;
        }

        .flow-item {
            padding: 0.375rem;
        }

        .flow-title {
            font-size: 0.6875rem;
        }
    }
</style>