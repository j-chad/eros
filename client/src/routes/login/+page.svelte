<script lang="ts">
	import {onDestroy, onMount} from 'svelte';
	import jsQR from 'jsqr';
	import BrandHeader from "$lib/ui/BrandHeader.svelte";
	import Button from "$lib/ui/base/Button.svelte";
	import {ArrowLeft, ArrowRight} from "lucide-svelte";
	import Card from "$lib/ui/base/Card.svelte";
	import {login} from "$lib/services/auth";

	type Mode = 'scan' | 'manual';
	let mode = $state<Mode>('scan');
	let step = $state<1 | 2>(1);

	let regCodeRaw = $state('');
	const regCode = $derived(formatRegCode(regCodeRaw));
	const regCodeValid = $derived(isValidRegCode(regCode));

	let deviceName = $state('');
	const deviceNameValid = $derived(deviceName.trim().length >= 2);

	let busy = $state(false);
	let error = $state<string | null>(null);

	// camera / scanning
	let canScan = $state(false);
	let scanning = $state(false);

	let stream: MediaStream | null = null;
	let rafId: number | null = null;

	let videoEl = $state<HTMLVideoElement | null>(null);
	let canvasEl = $state<HTMLCanvasElement | null>(null);
	let codeEl = $state<HTMLInputElement | null>(null);
	let nameEl = $state<HTMLInputElement | null>(null);

	onMount(async () => {
		canScan =
			typeof window !== 'undefined' &&
			!!navigator.mediaDevices?.getUserMedia &&
			!!window.isSecureContext; // camera requires https or localhost

		if (!canScan) {
			mode = 'manual';
			queueMicrotask(() => codeEl?.focus());
			return;
		}

		// QR scan is default
		await startScan();
	});

	onDestroy(() => stopScan());

	function formatRegCode(value: string): string {
		return value
	}

	function isValidRegCode(value: string) {
		return true
	}

	function parseRegisterUrlToCode(raw: string): string | null {
		// QR content should be a URL that ends with /register?code=...
		// Accept absolute or relative URLs (relative resolved to current origin).
		let url: URL;
		try {
			url = new URL(raw, window.location.origin);
		} catch {
			return null;
		}

		// Ensure path ends with /register (allow trailing slash)
		const path = url.pathname.replace(/\/+$/, '');
		if (!path.endsWith('/register')) return null;

		const codeParam = url.searchParams.get('code') ?? '';
		if (!codeParam) return null;

		// Convert code param into our expected XXXX-XXXX (accept 8 raw alphanum or dashed)
		const upper = codeParam.toUpperCase();
		// const dashed = upper.match(/^[A-Z0-9]{4}-[A-Z0-9]{4}$/)?.[0];
		// if (dashed) return dashed;
		//
		// const raw8 = upper.replace(/[^A-Z0-9]/g, '').match(/^[A-Z0-9]{8}$/)?.[0];
		// if (raw8) return `${raw8.slice(0, 4)}-${raw8.slice(4)}`;

		return null;
	}

	async function startScan() {
		error = null;
		if (!canScan || scanning) return;

		try {
			scanning = true;

			stream = await navigator.mediaDevices.getUserMedia({
				video: {
					facingMode: { ideal: 'environment' },
					width: { ideal: 1280 },
					height: { ideal: 720 }
				},
				audio: false
			});

			if (!videoEl) return;
			videoEl.srcObject = stream;
			videoEl.playsInline = true;
			videoEl.muted = true;
			await videoEl.play();

			const tick = () => {
				if (!scanning || !videoEl || !canvasEl) return;

				const w = videoEl.videoWidth;
				const h = videoEl.videoHeight;

				if (w && h) {
					// Keep canvas small-ish for mobile perf (downscale)
					const targetW = Math.min(720, w);
					const scale = targetW / w;
					const targetH = Math.round(h * scale);

					canvasEl.width = targetW;
					canvasEl.height = targetH;

					const ctx = canvasEl.getContext('2d', { willReadFrequently: true });
					if (ctx) {
						ctx.drawImage(videoEl, 0, 0, targetW, targetH);
						const imageData = ctx.getImageData(0, 0, targetW, targetH);

						const qr = jsQR(imageData.data, imageData.width, imageData.height, {
							inversionAttempts: 'attemptBoth'
						});

						if (qr?.data) {
							const parsed = parseRegisterUrlToCode(qr.data);
							if (parsed) {
								regCodeRaw = parsed;
								stopScan();
								goNext();
								return;
							}
						}
					}
				}

				rafId = requestAnimationFrame(tick);
			};

			rafId = requestAnimationFrame(tick);
		} catch {
			error = 'Couldn’t access the camera. Check permissions, or use manual entry.';
			stopScan();
			mode = 'manual';
			queueMicrotask(() => codeEl?.focus());
		}
	}

	function stopScan() {
		scanning = false;
		if (rafId != null) cancelAnimationFrame(rafId);
		rafId = null;

		if (stream) {
			for (const t of stream.getTracks()) t.stop();
		}
		stream = null;

		if (videoEl) videoEl.srcObject = null;
	}

	function switchToManual() {
		mode = 'manual';
		stopScan();
		queueMicrotask(() => codeEl?.focus());
	}

	async function switchToScan() {
		if (!canScan) return;
		mode = 'scan';
		await startScan();
	}

	function goNext() {
		error = null;
		if (!regCodeValid) return;

		step = 2;
		queueMicrotask(() => nameEl?.focus());
	}

	function back() {
		error = null;
		step = 1;
		if (mode === 'scan') startScan();
		else queueMicrotask(() => codeEl?.focus());
	}

	async function finish() {
		error = null;
		if (!deviceNameValid) return;

		busy = true;
		try {
			await login(regCode, deviceName);
			window.location.href = new URLSearchParams(window.location.search).get('returnTo') ?? '/';
		} catch (err) {
			console.error('Login failed', err);
			error = 'Unable to log in. Please try again.';
		} finally {
			busy = false;
		}
	}
</script>

<svelte:head>
	<title>Login • Eros</title>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<div class="min-h-dvh bg-linear-to-br from-pink-50 via-base-100 to-pink-100">
	<div class="mx-auto min-h-dvh max-w-md px-4 py-6">
		<BrandHeader subtitle="Link this device">
			{#snippet rightContent()}
				<ul class="steps steps-horizontal w-full max-w-45">
					<li class="step {step >= 1 ? 'step-secondary' : ''}">Code</li>
					<li class="step {step === 2 ? 'step-secondary' : ''}">Device</li>
				</ul>
			{/snippet}
		</BrandHeader>

		<Card padded class="mt-5">
			<div class="card-body gap-4">
				{#if error}
					<div class="alert alert-error rounded-2xl">
						<span class="text-sm">{error}</span>
					</div>
				{/if}

				{#if step === 1}
					{#if mode === 'scan'}
						<div class="relative overflow-hidden rounded-3xl bg-base-200 animate-popIn">
							<video
								bind:this={videoEl}
								class="aspect-3/4 w-full object-cover"
								playsinline
								muted
							></video>

							<!-- hidden canvas for decoding -->
							<canvas bind:this={canvasEl} class="hidden"></canvas>

							<div class="pointer-events-none absolute inset-0 grid place-items-center">
								<div class="w-[72%] max-w-70 aspect-square rounded-3xl border-2 border-white/70 shadow-[0_0_0_999px_rgba(0,0,0,0.25)]">
									<div class="h-full w-full rounded-3xl ring-2 ring-pink-300/60"></div>
								</div>
							</div>

							<div class="absolute bottom-3 left-3 right-3">
								<div class="rounded-2xl bg-base-100/80 backdrop-blur p-3 text-sm">
									{#if scanning}
										<div class="flex items-center gap-2">
											<span class="loading loading-spinner loading-sm"></span>
											Scanning…
										</div>
									{:else}
										<div class="opacity-80">Camera ready.</div>
									{/if}
								</div>
							</div>
						</div>

						<div class="divider my-1 opacity-60">Trouble scanning?</div>
						<button class="btn btn-ghost w-full rounded-2xl" onclick={switchToManual}>
							Enter code manually
						</button>

					{:else}
						<label class="form-control w-full animate-popIn">
							<div class="label">
								<span class="label-text font-semibold">Registration code</span>
								<span class="label-text-alt opacity-60">XXXX-XXXX-XXXX</span>
							</div>
							<input
								bind:this={codeEl}
								class="input input-bordered input-secondary w-full rounded-2xl text-lg tracking-widest"
								inputmode="text"
								autocomplete="one-time-code"
								placeholder="ABCD-1234-XYZ9"
								value={regCodeRaw}
								oninput={(e) => regCodeRaw = e.currentTarget.value}
								onkeydown={(e) => e.key === 'Enter' && goNext()}
							/>
							<div class="label">
								<span class="label-text-alt opacity-70">{regCodeValid ? 'Looks good ✨' : ' '}</span>
							</div>
						</label>

						<Button
							block
							disabled={!regCodeValid}
							onclick={goNext}
						>
							Continue
							<ArrowRight class="font-icon"/>
						</Button>

						{#if canScan}
							<div class="divider my-1 opacity-60">or</div>
							<Button block ghost onclick={switchToScan}>
								Back to scanning
							</Button>
						{/if}
					{/if}

				{:else}
					<div class="space-y-1 animate-popIn">
						<h1 class="text-xl font-bold">Name this device</h1>
						<p class="text-sm opacity-70">You can change it later.</p>
					</div>

					<div class="rounded-2xl bg-base-200/60 p-4 animate-popIn">
						<div class="text-xs font-semibold opacity-60">Registration code</div>
						<div class="font-mono text-sm">{regCode}</div>
					</div>

					<label class="form-control w-full animate-popIn">
						<div class="label">
							<span class="label-text font-semibold">Device name</span>
							<span class="label-text-alt opacity-60">2+ chars</span>
						</div>
						<input
							bind:this={nameEl}
							class="input input-bordered input-secondary w-full rounded-2xl text-lg"
							autocomplete="nickname"
							placeholder="e.g. Eros iPhone"
							bind:value={deviceName}
							onkeydown={(e) => e.key === 'Enter' && finish()}
						/>
						<div class="label">
              <span class="label-text-alt opacity-70">
                {deviceNameValid ? 'Nice ✨' : 'Give it a short name.'}
              </span>
						</div>
					</label>

					<div class="flex gap-3 pt-1 animate-popIn">
						<Button ghost disabled={busy} onclick={back}>
							<ArrowLeft class="font-icon"/>
							Back
						</Button>
						<button
							class="btn btn-secondary flex-1 rounded-2xl"
							disabled={!deviceNameValid || busy}
							onclick={finish}
						>
							{#if busy}
								<span class="loading loading-spinner loading-sm"></span>
								Saving…
							{:else}
								Finish
							{/if}
						</button>
					</div>
				{/if}
			</div>
		</Card>

		<div class="mt-5 text-center text-xs opacity-60">
			Tip: install the app for the smoothest scanning experience.
		</div>
	</div>
</div>

<style>
	:global(.animate-popIn) {
		animation: popIn 220ms ease-out both;
	}
	@keyframes popIn {
		from {
			opacity: 0;
			transform: translateY(6px) scale(0.99);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}
</style>
