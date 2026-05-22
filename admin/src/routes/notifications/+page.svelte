<script lang="ts">
    import { Send, Trash2, ChevronDown, ChevronUp } from 'lucide-svelte';
    import { api } from '$lib/api';
    import Card from '$lib/components/Card.svelte';
    import Button from '$lib/components/Button.svelte';
    import Header from '$lib/components/Header.svelte';
    import Table from '$lib/components/Table.svelte';
    import DateDisplay from '$lib/components/DateDisplay.svelte';
    import type { PushSubscription, PushSendRequest } from '$lib/types';

    let { data }: { data: { subscriptions: PushSubscription[] } } = $props();
    let subscriptions = $state(data.subscriptions);

    // --- Send form state ---
    let title = $state('');
    let body = $state('');
    let url = $state('');
    let tag = $state('');
    let topic = $state('');
    let urgency = $state<PushSendRequest['urgency']>('normal');
    let ttl = $state('');
    let showAdvanced = $state(false);
    let sending = $state(false);
    let result = $state<{ sent: number; failed: number; cleaned: number } | null>(null);

    function resetForm() {
        title = '';
        body = '';
        url = '';
        tag = '';
        topic = '';
        urgency = 'normal';
        ttl = '';
        result = null;
    }

    async function handleSend() {
        if (!title.trim()) {
            alert('Title is required');
            return;
        }
        if (!body.trim()) {
            alert('Body is required');
            return;
        }

        sending = true;
        result = null;

        try {
            const payload: PushSendRequest = {
                message: {
                    title: title.trim(),
                    body: body.trim(),
                    ...(tag.trim() && { tag: tag.trim() }),
                    ...(url.trim() && { data: { url: url.trim() } }),
                },
                ...(urgency && { urgency }),
                ...(topic.trim() && { topic: topic.trim() }),
                ...(ttl.trim() && { ttl: ttl.trim() }),
            };

            const sendResult = await api.push.send(payload);
            result = sendResult;
            title = '';
            body = '';
            url = '';
            tag = '';
            topic = '';
        } catch (err) {
            console.error('Failed to send notification:', err);
            alert('Failed to send notification');
        } finally {
            sending = false;
        }
    }

    async function handleDeleteSubscription(deviceId: string) {
        if (!confirm('Remove this push subscription? The device will need to re-enable notifications.')) return;

        try {
            await api.push.deleteSubscription(deviceId);
            subscriptions = subscriptions.filter(s => s.deviceID !== deviceId);
        } catch (err) {
            console.error('Failed to delete subscription:', err);
            alert('Failed to delete subscription');
        }
    }
</script>

<Header title="Notifications" />

<div class="page-content">
    <div class="grid">
        <Card title="Send Notification">
            <div class="form-group">
                <label for="title">Title <span class="required">*</span></label>
                <input
                    id="title"
                    type="text"
                    bind:value={title}
                    placeholder="e.g., New Adventure Available"
                    disabled={sending}
                />
            </div>

            <div class="form-group">
                <label for="body">Body <span class="required">*</span></label>
                <textarea
                    id="body"
                    bind:value={body}
                    placeholder="e.g., A new adventure is waiting for you"
                    rows="3"
                    disabled={sending}
                ></textarea>
            </div>

            <div class="form-group">
                <label for="url">Link</label>
                <input
                    id="url"
                    type="text"
                    bind:value={url}
                    placeholder="/ (default)"
                    disabled={sending}
                />
                <span class="hint">URL path to open when tapped. Leave empty for home screen.</span>
            </div>

            <button class="advanced-toggle" onclick={() => showAdvanced = !showAdvanced}>
                {#if showAdvanced}
                    <ChevronUp size={14} />
                {:else}
                    <ChevronDown size={14} />
                {/if}
                Advanced options
            </button>

            {#if showAdvanced}
                <div class="advanced-section">
                    <div class="grid-2">
                        <div class="form-group">
                            <label for="urgency">Urgency</label>
                            <select id="urgency" bind:value={urgency} disabled={sending}>
                                <option value="very-low">Very Low</option>
                                <option value="low">Low</option>
                                <option value="normal">Normal</option>
                                <option value="high">High</option>
                            </select>
                            <span class="hint">Affects delivery on battery-powered devices.</span>
                        </div>

                        <div class="form-group">
                            <label for="ttl">TTL</label>
                            <input
                                id="ttl"
                                type="text"
                                bind:value={ttl}
                                placeholder="e.g., 1h, 30m, 86400s"
                                disabled={sending}
                            />
                            <span class="hint">How long the push service retains the message if undelivered.</span>
                        </div>
                    </div>

                    <div class="grid-2">
                        <div class="form-group">
                            <label for="tag">Tag</label>
                            <input
                                id="tag"
                                type="text"
                                bind:value={tag}
                                placeholder="e.g., graph-unlock"
                                disabled={sending}
                            />
                            <span class="hint">Notifications with the same tag replace each other on the device.</span>
                        </div>

                        <div class="form-group">
                            <label for="topic">Topic</label>
                            <input
                                id="topic"
                                type="text"
                                bind:value={topic}
                                placeholder="e.g., graph-unlock"
                                disabled={sending}
                            />
                            <span class="hint">Messages with the same topic replace each other at the push service before delivery.</span>
                        </div>
                    </div>
                </div>
            {/if}

            {#if result}
                <div class="result">
                    Sent to {result.sent} {result.sent === 1 ? 'device' : 'devices'}{#if result.failed > 0}, {result.failed} failed{/if}{#if result.cleaned > 0}, {result.cleaned} expired subscriptions removed{/if}.
                </div>
            {/if}

            <div class="form-actions">
                <Button variant="primary" onclick={handleSend} disabled={sending}>
                    {#snippet icon()}<Send size={16} />{/snippet}
                    {sending ? 'Sending...' : 'Send'}
                </Button>
                {#if title || body || url || tag || topic}
                    <Button variant="secondary" onclick={resetForm} disabled={sending}>
                        Clear
                    </Button>
                {/if}
            </div>
        </Card>

        <Card title="Subscribed Devices">
            <Table headers={['Device', 'Subscribed', 'Last Updated', 'Expires', '']} empty="No devices subscribed to push notifications.">
                {#each subscriptions as sub (sub.deviceID)}
                    <tr>
                        <td>
                            <span class="device-name">{sub.device_name}</span>
                            <span class="device-endpoint" title={sub.endpoint}>
                                {new URL(sub.endpoint).hostname}
                            </span>
                        </td>
                        <td><DateDisplay datetime={sub.created_at} /></td>
                        <td><DateDisplay datetime={sub.updated_at} /></td>
                        <td>
                            {#if sub.expirationTime}
                                <DateDisplay datetime={sub.expirationTime} expiry />
                            {:else}
                                <span class="no-expiry">None</span>
                            {/if}
                        </td>
                        <td class="actions">
                            <Button variant="danger" size="sm" onclick={() => handleDeleteSubscription(sub.deviceID)}>
                                {#snippet icon()}<Trash2 size={14} />{/snippet}
                                Remove
                            </Button>
                        </td>
                    </tr>
                {/each}
            </Table>
        </Card>
    </div>
</div>

<style>
    .page-content {
        width: 100%;
    }

    .grid {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: #4b5563;
        margin-bottom: 0.5rem;
    }

    .required {
        color: #ef4444;
    }

    input[type="text"],
    textarea,
    select {
        width: 100%;
        padding: 0.5rem 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 4px;
        font-size: 0.875rem;
        font-family: inherit;
    }

    input:focus,
    textarea:focus,
    select:focus {
        outline: none;
        border-color: #3b82f6;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    input:disabled,
    textarea:disabled,
    select:disabled {
        background: #f3f4f6;
        cursor: not-allowed;
    }

    .hint {
        display: block;
        font-size: 0.75rem;
        color: #9ca3af;
        margin-top: 0.25rem;
    }

    .advanced-toggle {
        display: flex;
        align-items: center;
        gap: 0.25rem;
        background: none;
        border: none;
        color: #6b7280;
        font-size: 0.8125rem;
        cursor: pointer;
        padding: 0;
        margin-bottom: 1rem;
    }

    .advanced-toggle:hover {
        color: #374151;
    }

    .advanced-section {
        padding: 1rem;
        background: #f9fafb;
        border-radius: 6px;
        margin-bottom: 1rem;
    }

    .advanced-section .form-group {
        margin-bottom: 1rem;
    }

    .advanced-section .form-group:last-child {
        margin-bottom: 0;
    }

    .grid-2 {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }

    .form-actions {
        display: flex;
        gap: 0.5rem;
        margin-top: 1.5rem;
    }

    .result {
        background: #ecfdf5;
        color: #065f46;
        border: 1px solid #a7f3d0;
        border-radius: 4px;
        padding: 0.75rem;
        font-size: 0.875rem;
    }

    .device-name {
        display: block;
        font-weight: 500;
    }

    .device-endpoint {
        display: block;
        font-size: 0.75rem;
        color: #9ca3af;
    }

    .no-expiry {
        color: #9ca3af;
        font-size: 0.875rem;
    }

    .actions {
        text-align: right;
    }
</style>
