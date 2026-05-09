<script lang="ts">
	import { SvelteMap } from 'svelte/reactivity';
	import { ChevronLeft, ChevronRight, Pointer } from 'lucide-svelte';
	import type { GraphSummary } from '$lib/types/graph';
	import { KVStore, KVKey } from '$lib/db/stores/kv';

	const {
		graphs,
		onSelectDate,
	}: {
		graphs: GraphSummary[];
		onSelectDate?: (date: Date, graphs: GraphSummary[]) => void;
	} = $props();

	const WEEKDAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

	// --- helpers ---

	/** Format a Date to a local YYYY-MM-DD key for grouping. */
	function toDateKey(d: Date): string {
		const y = d.getFullYear();
		const m = String(d.getMonth() + 1).padStart(2, '0');
		const day = String(d.getDate()).padStart(2, '0');
		return `${y}-${m}-${day}`;
	}

	function sameMonth(a: Date, b: Date): boolean {
		return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth();
	}

	// --- state ---

	let currentMonth = $state(new Date());

	// null = not yet loaded from KV (suppress cursor until known to avoid flicker)
	let showCursor = $state<boolean | null>(null);

	$effect(() => {
		KVStore.get(KVKey.CalendarTipSeen).then((seen) => {
			showCursor = !seen;
		});
	});

	// --- derived ---

	const now = new Date();
	const todayKey = toDateKey(now);

	/** Map of YYYY-MM-DD → GraphSummary[] for O(1) lookup. */
	const graphsByDay = $derived.by(() => {
		const map = new SvelteMap<string, GraphSummary[]>();
		for (const g of graphs) {
			const key = toDateKey(g.starting_at);
			const arr = map.get(key);
			if (arr) {
				arr.push(g);
			} else {
				map.set(key, [g]);
			}
		}
		return map;
	});

	/** Sorted list of all graph dates for determining navigation bounds. */
	const sortedDates = $derived(
		[...graphs].sort((a, b) => a.starting_at.getTime() - b.starting_at.getTime()),
	);

	const minMonth = $derived(sortedDates.length > 0 ? sortedDates[0].starting_at : now);
	const maxMonth = $derived(
		sortedDates.length > 0 ? sortedDates[sortedDates.length - 1].starting_at : now,
	);

	const canGoBack = $derived(!sameMonth(currentMonth, minMonth));
	const canGoForward = $derived(!sameMonth(currentMonth, maxMonth));

	const monthLabel = $derived(
		currentMonth.toLocaleDateString(undefined, { month: 'long', year: 'numeric' }),
	);

	/** Build the grid of day cells for the current month. */
	const dayCells = $derived.by(() => {
		const year = currentMonth.getFullYear();
		const month = currentMonth.getMonth();

		const firstOfMonth = new Date(year, month, 1);
		let startWeekday = firstOfMonth.getDay() - 1; // JS Sunday=0 → Mon=-1
		if (startWeekday < 0) startWeekday = 6; // Sunday → 6

		const daysInMonth = new Date(year, month + 1, 0).getDate();

		const cells: ({ day: number; key: string } | null)[] = [];

		for (let i = 0; i < startWeekday; i++) {
			cells.push(null);
		}

		for (let d = 1; d <= daysInMonth; d++) {
			cells.push({ day: d, key: toDateKey(new Date(year, month, d)) });
		}

		while (cells.length % 7 !== 0) {
			cells.push(null);
		}

		return cells;
	});

	/** Key of the first past-graph day in the visible month - this is where the cursor sits. */
	const cursorKey = $derived.by(() => {
		for (const cell of dayCells) {
			if (!cell) continue;
			const dayGraphs = graphsByDay.get(cell.key);
			if (dayGraphs?.some((g) => g.starting_at.getTime() <= Date.now())) {
				return cell.key;
			}
		}
		return null;
	});

	// --- actions ---

	function goBack() {
		if (!canGoBack) return;
		currentMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1, 1);
	}

	function goForward() {
		if (!canGoForward) return;
		currentMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 1);
	}

	async function handleDayClick(key: string) {
		const dayGraphs = graphsByDay.get(key);
		if (!dayGraphs?.length) return;

		// First tap on any past day dismisses the cursor hint.
		if (showCursor) {
			showCursor = false;
			await KVStore.set(KVKey.CalendarTipSeen, true);
		}

		if (!onSelectDate) return;
		const [year, month, day] = key.split('-').map(Number);
		onSelectDate(new Date(year, month - 1, day), dayGraphs);
	}
</script>

<div class="flex flex-col gap-3">
	<!-- Month navigation header -->
	<div class="flex items-center justify-between px-1">
		<button
			class="btn btn-ghost btn-sm btn-circle"
			disabled={!canGoBack}
			onclick={goBack}
			aria-label="Previous month"
		>
			<ChevronLeft size={18} />
		</button>
		<span class="text-base font-bold">{monthLabel}</span>
		<button
			class="btn btn-ghost btn-sm btn-circle"
			disabled={!canGoForward}
			onclick={goForward}
			aria-label="Next month"
		>
			<ChevronRight size={18} />
		</button>
	</div>

	<!-- Weekday labels -->
	<div class="grid grid-cols-7 text-center">
		{#each WEEKDAYS as day (day)}
			<span class="text-xs font-semibold opacity-50 py-1">{day}</span>
		{/each}
	</div>

	<!-- Day grid -->
	<div class="grid grid-cols-7">
		{#each dayCells as cell, i (cell?.key ?? `blank-${i}`)}
		{@const dayGraphs = cell ? graphsByDay.get(cell.key) : undefined}
		{@const hasGraphs = !!dayGraphs?.length}
		{@const isToday = cell?.key === todayKey}
		{@const isPast = hasGraphs && dayGraphs.some((g) => g.starting_at.getTime() <= Date.now())}
		{@const isFuture = hasGraphs && !isPast}
		{@const hasCursor = showCursor === true && cell?.key === cursorKey}
		{@const dotClass = isPast
			? dayGraphs?.some((g) => g.status === 'completed')
				? 'bg-success'
				: dayGraphs?.some((g) => g.status === 'in_progress')
					? 'bg-warning'
					: 'bg-primary'
			: 'bg-base-300'}

			{#if cell === null}
				<div></div>
			{:else}
				<button
					class="relative overflow-visible flex flex-col items-center justify-center gap-0.5 py-3 rounded-xl transition-colors
						{isToday ? 'ring-2 ring-primary/30 ring-inset' : ''}
						{isPast ? 'cursor-pointer hover:bg-primary/10' : ''}
						{!hasGraphs || isFuture ? 'cursor-default' : ''}"
					onclick={() => isPast && handleDayClick(cell.key)}
					disabled={!isPast}
					aria-label="{cell.day}{hasGraphs ? ', has activity' : ''}{isToday ? ', today' : ''}"
				>
					<span
						class="text-sm leading-none
							{isToday ? 'font-bold text-primary' : ''}
							{isPast ? 'font-semibold' : ''}
							{!hasGraphs && !isToday ? 'opacity-60' : ''}"
					>
						{cell.day}
					</span>
				{#if isPast || isFuture}
					<span class="w-1.5 h-1.5 rounded-full {dotClass} {isToday ? 'animate-softPulse' : ''}"></span>
				{:else}
					<span class="w-1.5 h-1.5"></span>
				{/if}

					{#if hasCursor}
						<span class="absolute -bottom-0.5 left-1/2 -rotate-30 text-primary animate-cursorTap pointer-events-none">
							<Pointer size={16} />
						</span>
					{/if}
				</button>
			{/if}
		{/each}
	</div>
</div>
