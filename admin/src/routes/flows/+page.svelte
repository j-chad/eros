<!-- src/routes/flows/+page.svelte -->
<script lang="ts">
    import { ChevronLeft, ChevronRight, Plus, Edit, Trash2 } from 'lucide-svelte';
    import Button from '$lib/components/Button.svelte';
    import Card from '$lib/components/Card.svelte';
    import Header from '$lib/components/Header.svelte';
    import { api } from '$lib/api';

    type FlowStart = {
        id: number;
        title: string;
        description: string | null;
        startingAt: string;
    };

    let { data } = $props();

    let flows = $state<FlowStart[]>(data.flows);
    let currentDate = $state(new Date());

    $effect(() => {
        flows = data.flows;
    });

    // Calendar calculations
    let year = $derived(currentDate.getFullYear());
    let month = $derived(currentDate.getMonth());

    let monthName = $derived(
        new Date(year, month).toLocaleString('default', { month: 'long', year: 'numeric' })
    );

    let daysInMonth = $derived(new Date(year, month + 1, 0).getDate());
    let firstDayOfMonth = $derived(new Date(year, month, 1).getDay());

    let calendarDays = $derived(() => {
        const days = [];
        const prevMonthDays = new Date(year, month, 0).getDate();

        // Previous month days
        for (let i = firstDayOfMonth - 1; i >= 0; i--) {
            days.push({
                day: prevMonthDays - i,
                isCurrentMonth: false,
                date: new Date(year, month - 1, prevMonthDays - i)
            });
        }

        // Current month days
        for (let i = 1; i <= daysInMonth; i++) {
            days.push({
                day: i,
                isCurrentMonth: true,
                date: new Date(year, month, i)
            });
        }

        // Next month days to fill the grid
        const remainingDays = 42 - days.length; // 6 rows × 7 days
        for (let i = 1; i <= remainingDays; i++) {
            days.push({
                day: i,
                isCurrentMonth: false,
                date: new Date(year, month + 1, i)
            });
        }

        return days;
    });

    function getFlowsForDate(date: Date): FlowStart[] {
        const dateStr = date.toISOString().split('T')[0];
        return flows.filter(flow => {
            const flowDate = new Date(flow.startingAt).toISOString().split('T')[0];
            return flowDate === dateStr;
        });
    }

    function previousMonth() {
        currentDate = new Date(year, month - 1);
    }

    function nextMonth() {
        currentDate = new Date(year, month + 1);
    }

    function today() {
        currentDate = new Date();
    }

    function isToday(date: Date): boolean {
        const today = new Date();
        return date.toDateString() === today.toDateString();
    }

    async function handleCreateFlow(date: Date) {
        // Navigate to create flow page with date
        console.log('Create flow for:', date);
        // window.location.href = `/flows/create?date=${date.toISOString()}`;
    }

    async function handleEditFlow(flowId: number) {
        console.log('Edit flow:', flowId);
        // window.location.href = `/flows/${flowId}/edit`;
    }

    async function handleDeleteFlow(flowId: number) {
        if (confirm('Are you sure you want to delete this flow?')) {
            try {
                await api.flows.delete(flowId);
                flows = flows.filter(f => f.id !== flowId);
            } catch (error) {
                console.error('Failed to delete flow:', error);
            }
        }
    }
</script>

<Header title="Reward Flows Calendar" />

<Card title={monthName}>
    {#snippet actions()}
        <div class="calendar-controls">
            <Button variant="secondary" size="sm" onclick={today}>
                Today
            </Button>
            <div class="month-nav">
                <button class="nav-btn" onclick={previousMonth}>
                    <ChevronLeft size={20} />
                </button>
                <button class="nav-btn" onclick={nextMonth}>
                    <ChevronRight size={20} />
                </button>
            </div>
        </div>
    {/snippet}

    <div class="calendar">
        <!-- Day headers -->
        <div class="calendar-header">
            {#each ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as day}
                <div class="day-header">{day}</div>
            {/each}
        </div>

        <!-- Calendar grid -->
        <div class="calendar-grid">
            {#each calendarDays() as { day, isCurrentMonth, date }}
                <div
                        class="calendar-day"
                        class:other-month={!isCurrentMonth}
                        class:today={isToday(date)}
                >
                    <div class="day-number">{day}</div>

                    {#if isCurrentMonth}
                        {@const dayFlows = getFlowsForDate(date)}

                        {#if dayFlows.length > 0}
                            <div class="flows">
                                {#each dayFlows as flow}
                                    <div class="flow-item">
                                        <div class="flow-content">
                                            <div class="flow-title">{flow.title}</div>
                                            {#if flow.description}
                                                <div class="flow-description">{flow.description}</div>
                                            {/if}
                                            <div class="flow-time">
                                                {new Date(flow.startingAt).toLocaleTimeString('en-US', {
                                                    hour: '2-digit',
                                                    minute: '2-digit'
                                                })}
                                            </div>
                                        </div>
                                        <div class="flow-actions">
                                            <button
                                                    class="action-btn edit"
                                                    onclick={() => handleEditFlow(flow.id)}
                                                    title="Edit flow"
                                            >
                                                <Edit size={12} />
                                            </button>
                                            <button
                                                    class="action-btn delete"
                                                    onclick={() => handleDeleteFlow(flow.id)}
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
                                onclick={() => handleCreateFlow(date)}
                                title="Create flow"
                        >
                            <Plus size={14} />
                        </button>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</Card>

<style>
    .calendar-controls {
        display: flex;
        gap: 1rem;
        align-items: center;
    }

    .month-nav {
        display: flex;
        gap: 0.25rem;
    }

    .nav-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2rem;
        height: 2rem;
        background: #f3f4f6;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        cursor: pointer;
        transition: background 0.2s;
    }

    .nav-btn:hover {
        background: #e5e7eb;
    }

    .calendar {
        width: 100%;
    }

    .calendar-header {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        gap: 1px;
        background: #e5e7eb;
        border: 1px solid #e5e7eb;
        border-bottom: none;
    }

    .day-header {
        background: #f9fafb;
        padding: 0.75rem;
        text-align: center;
        font-size: 0.75rem;
        font-weight: 600;
        color: #6b7280;
        text-transform: uppercase;
    }

    .calendar-grid {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        gap: 1px;
        background: #e5e7eb;
        border: 1px solid #e5e7eb;
    }

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
        .day-header {
            font-size: 0.625rem;
            padding: 0.5rem 0.25rem;
        }

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

        .flow-time {
            font-size: 0.625rem;
        }
    }
</style>