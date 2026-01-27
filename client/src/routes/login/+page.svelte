<script lang="ts">
	let registrationCode = "";
	let deviceName = "";
	let isSubmitting = false;

	function normalizeCode(v: string) {
		return v.toUpperCase().replace(/\s+/g, "");
	}

	async function onSubmit() {
		isSubmitting = true;
		try {
			// Wire this to a SvelteKit action or endpoint later
			await new Promise((r) => setTimeout(r, 500));
			alert("Activated (demo). Hook this to your backend/action.");
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Eros • Activate</title>
	<meta name="apple-mobile-web-app-capable" content="yes" />
	<meta name="theme-color" content="#ff4d8d" />
</svelte:head>

<main
	class="
    min-h-screen bg-base-200 text-base-content
    px-4 py-8
    flex items-center justify-center
  "
	style="
    background-image:
      radial-gradient(900px 500px at 15% 10%, oklch(var(--p) / 0.12), transparent 55%),
      radial-gradient(900px 500px at 85% 30%, oklch(var(--s) / 0.10), transparent 60%),
      radial-gradient(900px 500px at 50% 90%, oklch(var(--a) / 0.08), transparent 55%);
  "
>
	<div class="w-full max-w-md">
		<!-- Header -->
		<header class="mb-4 px-1">
			<div class="flex items-start gap-3">
				<div class="h-12 w-12 rounded-2xl bg-base-100 shadow grid place-items-center">
					<span class="text-xl">💗</span>
				</div>

				<div class="min-w-0">
					<h1 class="text-3xl font-bold leading-tight">
						Eros
					</h1>
					<p class="mt-1 text-sm opacity-80">
						Enter your <span class="font-semibold">registration code</span> and name this device.
					</p>
				</div>
			</div>
		</header>

		<!-- Card -->
		<section class="card bg-base-100 shadow-xl rounded-3xl">
			<div class="card-body gap-4">
				<form class="flex flex-col gap-3" on:submit|preventDefault={onSubmit}>
					<!-- Registration Code -->
					<div class="form-control">
						<label class="label" for="code">
							<span class="label-text">Registration code</span>
							<span class="label-text-alt opacity-70">Case-insensitive</span>
						</label>

						<label class="input input-bordered rounded-2xl flex items-center gap-2">
							<span class="opacity-70">🔑</span>
							<input
								id="code"
								name="code"
								class="grow"
								inputmode="latin"
								autocomplete="one-time-code"
								placeholder="e.g. LOVE-8K2Q"
								bind:value={registrationCode}
								on:input={(e) => (registrationCode = normalizeCode((e.target as HTMLInputElement).value))}
								required
							/>
						</label>

						<label class="label">
              <span class="label-text-alt opacity-70">
                Printed on the card you were given.
              </span>
						</label>
					</div>

					<!-- Device Name -->
					<div class="form-control">
						<label class="label" for="device">
							<span class="label-text">Device name</span>
							<span class="label-text-alt opacity-70">You can change this later</span>
						</label>

						<label class="input input-bordered rounded-2xl flex items-center gap-2">
							<span class="opacity-70">📱</span>
							<input
								id="device"
								name="device"
								class="grow"
								autocomplete="nickname"
								placeholder="e.g. Jackson’s Phone ✨"
								bind:value={deviceName}
								required
							/>
						</label>

						<label class="label">
              <span class="label-text-alt opacity-70">
                Helps Eros recognize this device.
              </span>
						</label>
					</div>

					<!-- Actions -->
					<div class="mt-2 flex flex-col gap-2">
						<button
							class="btn btn-primary rounded-2xl w-full"
							type="submit"
							disabled={isSubmitting || registrationCode.length < 4 || deviceName.trim().length < 2}
						>
							{#if isSubmitting}
								<span class="loading loading-spinner"></span>
								Activating…
							{:else}
								Continue 💞
							{/if}
						</button>

						<div class="text-xs sm:text-sm opacity-75 text-center">
							Tip: for the best experience, install <span class="font-semibold">Eros</span> from your browser menu.
						</div>
					</div>
				</form>

				<div class="divider my-2">Help</div>

				<div class="collapse collapse-arrow bg-base-200 rounded-2xl">
					<input type="checkbox" />
					<div class="collapse-title font-medium">
						Where do I find the registration code?
					</div>
					<div class="collapse-content">
						<p class="opacity-80 text-sm">
							It’s on the physical card/printout you received. Upper/lowercase doesn’t matter.
							If it includes hyphens, keep them.
						</p>
					</div>
				</div>

				<div class="mt-2 flex items-center gap-2">
					<span class="badge badge-secondary rounded-xl">Private</span>
					<span class="text-sm opacity-80">
            Your code is only used to activate Eros on this device.
          </span>
				</div>
			</div>
		</section>

		<footer class="mt-4 text-center">
			<p class="text-xs opacity-60">Made with love 🩷</p>
		</footer>
	</div>
</main>
