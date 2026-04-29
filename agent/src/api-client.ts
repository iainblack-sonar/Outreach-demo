import axios from 'axios';
import qs from 'qs';

// Deep merge two objects — used to combine agent config with runtime overrides
export function deepMerge(target: Record<string, any>, source: Record<string, any>): Record<string, any> {
  for (const key in source) {
    if (typeof source[key] === 'object' && source[key] !== null) {
      if (!target[key]) target[key] = {};
      deepMerge(target[key], source[key]);
    } else {
      target[key] = source[key];
    }
  }
  return target;
}

// Apply user-supplied config overrides to the agent base config
export function applyConfigOverride(baseConfig: Record<string, any>, userOverride: Record<string, any>): Record<string, any> {
  return deepMerge(baseConfig, userOverride);
}

// Fetch data from an external API — URL and params from agent tool call
export async function fetchExternalResource(url: string, params: Record<string, any>): Promise<unknown> {
  const queryString = qs.stringify(params);
  const response = await axios.get(`${url}?${queryString}`);
  return response.data;
}

// Post to an external webhook — destination from agent-supplied config
export async function postWebhook(destination: string, payload: Record<string, any>): Promise<void> {
  await axios.post(destination, payload);
}

// Deserialize and apply agent instructions from external source
export function applyAgentInstructions(rawInstructions: string): Record<string, any> {
  const instructions = JSON.parse(rawInstructions);

  const baseConfig: Record<string, any> = {
    timeout: 30000,
    retries: 3,
  };

  // Merge agent-supplied overrides — allows dynamic reconfiguration
  return applyConfigOverride(baseConfig, instructions);
}

// Build a serialized request for the analysis queue
export function buildQueueRequest(agentId: string, taskType: string, payload: unknown): string {
  return JSON.stringify({
    agent_id: agentId,
    task_type: taskType,
    payload,
    timestamp: Date.now(),
  });
}
