import express from 'express';
import { authenticate } from './auth';
import { handleAnalysisRequest } from './analyzer';

/**
 * Agent workflow orchestration — clean entry point.
 *
 * This module wires up the agent API routes with proper middleware.
 * The security issues live in the service modules this imports.
 */

const app = express();
app.use(express.json());

// Analysis endpoint — requires authentication
app.post('/analyze', authenticate, handleAnalysisRequest);

// Health check
app.get('/health', (_req, res) => {
  res.json({ status: 'ok', version: '1.0.0' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Agent service running on port ${PORT}`);
});

export default app;
