<script lang="ts">
	import { onMount } from 'svelte';

	const { duration = 3500 }: { duration?: number } = $props();

	let canvas: HTMLCanvasElement;

	// Pinks, roses, and warm accent colours that suit the valentine theme.
	const COLORS = [
		'#e8457e', // primary rose
		'#f472b6', // pink-400
		'#fb7185', // rose-400
		'#fda4af', // rose-300
		'#f9a8d4', // pink-300
		'#fbbf24', // amber-400 (warm sparkle)
		'#c084fc', // purple-400
		'#ffffff', // white
	];

	interface Particle {
		x: number;
		y: number;
		vx: number;
		vy: number;
		size: number;
		color: string;
		rotation: number;
		rotationSpeed: number;
		shape: 'circle' | 'rect' | 'heart';
		opacity: number;
		gravity: number;
	}

	function createParticles(cx: number, cy: number, count: number): Particle[] {
		const particles: Particle[] = [];
		for (let i = 0; i < count; i++) {
			const angle = Math.random() * Math.PI * 2;
			const speed = 1.5 + Math.random() * 4;
			const shapes: Particle['shape'][] = ['circle', 'rect', 'heart'];
			particles.push({
				x: cx,
				y: cy,
				vx: Math.cos(angle) * speed,
				vy: Math.sin(angle) * speed - 2.5, // gentle upward bias
				size: 5 + Math.random() * 7,
				color: COLORS[Math.floor(Math.random() * COLORS.length)],
				rotation: Math.random() * Math.PI * 2,
				rotationSpeed: (Math.random() - 0.5) * 0.1,
				shape: shapes[Math.floor(Math.random() * shapes.length)],
				opacity: 1,
				gravity: 0.04 + Math.random() * 0.03,
			});
		}
		return particles;
	}

	function drawHeart(ctx: CanvasRenderingContext2D, size: number) {
		const s = size * 0.5;
		ctx.beginPath();
		ctx.moveTo(0, s * 0.4);
		ctx.bezierCurveTo(-s, -s * 0.4, -s * 0.6, -s, 0, -s * 0.4);
		ctx.bezierCurveTo(s * 0.6, -s, s, -s * 0.4, 0, s * 0.4);
		ctx.fill();
	}

	onMount(() => {
		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		// fixed inset-0 already makes the canvas fill the viewport via CSS.
		// Set the buffer to match the CSS layout size (no DPR scaling needed
		// for a brief decorative animation — keeps coordinates simple).
		const rect = canvas.getBoundingClientRect();
		const w = rect.width;
		const h = rect.height;
		canvas.width = w;
		canvas.height = h;

		// Burst from the horizontal centre, roughly upper-third of screen
		const cx = w / 2;
		const cy = h * 0.38;
		const particles = createParticles(cx, cy, 80);

		const start = performance.now();
		let frame: number;

		function tick() {
			const elapsed = performance.now() - start;
			if (elapsed > duration) {
				ctx!.clearRect(0, 0, w, h);
				return;
			}

			ctx!.clearRect(0, 0, w, h);
			const progress = elapsed / duration;

			for (const p of particles) {
				p.x += p.vx;
				p.y += p.vy;
				p.vy += p.gravity;
				p.vx *= 0.98;
				p.rotation += p.rotationSpeed;

				// Fade out in the last 40% of the animation
				p.opacity = progress > 0.6 ? 1 - (progress - 0.6) / 0.4 : 1;

				ctx!.save();
				ctx!.translate(p.x, p.y);
				ctx!.rotate(p.rotation);
				ctx!.globalAlpha = p.opacity;
				ctx!.fillStyle = p.color;

				if (p.shape === 'circle') {
					ctx!.beginPath();
					ctx!.arc(0, 0, p.size / 2, 0, Math.PI * 2);
					ctx!.fill();
				} else if (p.shape === 'rect') {
					ctx!.fillRect(-p.size / 2, -p.size / 4, p.size, p.size / 2);
				} else {
					drawHeart(ctx!, p.size);
				}

				ctx!.restore();
			}

			frame = requestAnimationFrame(tick);
		}

		frame = requestAnimationFrame(tick);

		return () => cancelAnimationFrame(frame);
	});
</script>

<canvas
	bind:this={canvas}
	class="pointer-events-none fixed inset-0 z-50"
	aria-hidden="true"
></canvas>
