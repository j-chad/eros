<script lang="ts">
    import {ChevronLeft, ChevronRight} from 'lucide-svelte';
    import Button from '$lib/components/Button.svelte';
    import Card from '$lib/components/Card.svelte';
    import type {StartNode} from "$lib/types";
    import CalendarDay from "./CalendarDay.svelte";

    let {nodes, onDeleteGraph, onCreateGraph, onEditGraph}: {
        nodes: StartNode[];
        onCreateGraph?: (date: Date) => Promise<void>;
        onEditGraph?: (id: string) => Promise<void>;
        onDeleteGraph?: (id: string) => Promise<void>;
    } = $props();

    let currentDate = $state(new Date());

    // Calendar calculations
    let year = $derived(currentDate.getFullYear());
    let month = $derived(currentDate.getMonth());

    let monthName = $derived(
        new Date(year, month).toLocaleString('default', {month: 'long', year: 'numeric'})
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

    function getNodesForDate(date: Date): StartNode[] {
        const dateStr = date.toISOString().split('T')[0];
        return nodes.filter(node => {
            const graphDate = new Date(node.start.starting_at).toISOString().split('T')[0];
            return graphDate === dateStr;
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
</script>

<Card title={monthName}>
    {#snippet actions()}
        <div class="calendar-controls">
            <Button variant="secondary" size="sm" onclick={today}>
                Today
            </Button>
            <div class="month-nav">
                <button class="nav-btn" onclick={previousMonth}>
                    <ChevronLeft size={20}/>
                </button>
                <button class="nav-btn" onclick={nextMonth}>
                    <ChevronRight size={20}/>
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
            {#each calendarDays() as {day, isCurrentMonth, date}}
                <CalendarDay nodes={getNodesForDate(date)} day={day} date={date}
                             isCurrentMonth={isCurrentMonth} onCreateGraph={onCreateGraph}
                             onDeleteGraph={onDeleteGraph} onEditGraph={onEditGraph}></CalendarDay>
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

    /*.calendar {*/
    /*    width: 100%;*/
    /*}*/

    /*.calendar-header {*/
    /*    display: grid;*/
    /*    grid-template-columns: repeat(7, 1fr);*/
    /*    gap: 1px;*/
    /*    background: #e5e7eb;*/
    /*    border: 1px solid #e5e7eb;*/
    /*    border-bottom: none;*/
    /*}*/

    /*.day-header {*/
    /*    background: #f9fafb;*/
    /*    padding: 0.75rem;*/
    /*    text-align: center;*/
    /*    font-size: 0.75rem;*/
    /*    font-weight: 600;*/
    /*    color: #6b7280;*/
    /*    text-transform: uppercase;*/
    /*}*/

    /*.calendar-grid {*/
    /*    display: grid;*/
    /*    grid-template-columns: repeat(7, 1fr);*/
    /*    gap: 1px;*/
    /*    background: #e5e7eb;*/
    /*    border: 1px solid #e5e7eb;*/
    /*}*/

    .calendar {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        gap: 1px;
        background: #e5e7eb;
        border: 1px solid #e5e7eb;
        width: 100%;
    }

    .calendar-header {
        display: contents;
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
        display: contents;
    }

    @media (max-width: 768px) {
        .day-header {
            font-size: 0.625rem;
            padding: 0.5rem 0.25rem;
        }
    }
</style>