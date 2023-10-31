package llamaservergo

type PropsGetResponse struct {
	AssistantName string `json:"assistant_name"`
	AntiPrompt    string `json:"anti_prompt"`
	UserName      string `json:"user_name"`
}

type CompletionPostRequest struct {
	CompletionSettings
	Prompt string `json:"prompt"`
}

type CompletionSettings struct {
	// Temperature adjusts the randomness of the generated text (default: 0.8).
	Temperature *float64 `json:"temperature"`
	// TopK limits the next token selection to the K most probable tokens (default: 40).
	TopK *int `json:"top_k"`
	// TopP limits the next token selection to a subset of tokens with a cumulative probability above a threshold P (default: 0.95).
	TopP *float64 `json:"top_p"`
	// NPredict sets the maximum number of tokens to predict when generating text. **Note:** May exceed the set limit slightly if the last token is a partial multibyte character. When 0, no tokens will be generated but the prompt is evaluated into the cache. (default: -1, -1 = infinity).
	NPredict *int `json:"n_predict"`
	// NKeep specifies the number of tokens from the prompt to retain when the context size is exceeded and tokens need to be discarded.
	// By default, this value is set to 0 (meaning no tokens are kept). Use `-1` to retain all tokens from the prompt.
	NKeep *int `json:"n_keep"`
	// Stream allows receiving each predicted token in real-time instead of waiting for the completion to finish. To enable this, set to `true`.
	Stream bool `json:"stream"`
	// Stop specifies a JSON array of stopping strings.
	// These words will not be included in the completion, so make sure to add them to the prompt for the next iteration (default: []).
	Stop []string `json:"stop"`
	// TfsZ enables tail free sampling with parameter z (default: 1.0, 1.0 = disabled).
	TfsZ *float64 `json:"tfs_z"`
	// TypicalP enables locally typical sampling with parameter p (default: 1.0, 1.0 = disabled).
	TypicalP *float64 `json:"typical_p"`
	// RepeatPenalty controls the repetition of token sequences in the generated text (default: 1.1).
	RepeatPenalty *float64 `json:"repeat_penalty"`
	// RepeatLastN specifies the last n tokens to consider for penalizing repetition (default: 64, 0 = disabled, -1 = ctx-size).
	RepeatLastN *int `json:"repeat_last_n"`
	// PenalizeNl penalizes newline tokens when applying the repeat penalty (default: true).
	PenalizeNl *bool `json:"penalize_nl"`
	// PresencePenalty specifies the repeat alpha presence penalty (default: 0.0, 0.0 = disabled).
	PresencePenalty float64 `json:"presence_penalty"`
	// FrequencyPenalty specifies the repeat alpha frequency penalty (default: 0.0, 0.0 = disabled);
	FrequencyPenalty float64 `json:"frequency_penalty"`
	// Mirostat enables Mirostat sampling, controlling perplexity during text generation (default: 0, 0 = disabled, 1 = Mirostat, 2 = Mirostat 2.0).
	Mirostat int `json:"mirostat"`
	// MirostatTau specifies the Mirostat target entropy, parameter tau (default: 5.0).
	MirostatTau *float64 `json:"mirostat_tau"`
	// MirostatEta specifies the Mirostat learning rate, parameter eta (default: 0.1).
	MirostatEta *float64 `json:"mirostat_eta"`
	// Grammar specifies the grammar for grammar-based sampling (default: no grammar)
	Grammar *string `json:"grammar"`
	// Seed specifies the random number generator (RNG) seed (default: -1, -1 = random seed).
	Seed *int `json:"seed"`
	// IgnoreEos ignores end of stream token and continue generating (default: false).
	IgnoreEos *bool `json:"ignore_eos"`
	// LogitBias modifies the likelihood of a token appearing in the generated text completion. For example, use `"logit_bias": [[15043,1.0]]` to increase the likelihood of the token 'Hello', or `"logit_bias": [[15043,-1.0]]` to decrease its likelihood. Setting the value to false, `"logit_bias": [[15043,false]]` ensures that the token `Hello` is never produced (default: []).
	LogitBias [][]float64 `json:"logit_bias"`
	// NProbs specifies the number of top tokens to return with their probabilities (default: 0, 0 = disabled).
	NProbs int `json:"n_probs"`
	// ImageData specifies an array of objects to hold base64-encoded image `data` and its `id`s to be reference in `prompt`. You can determine the place of the image in the prompt as in the following: `USER:[img-12]Describe the image in detail.\nASSISTANT:` In this case, `[img-12]` will be replaced by the embeddings of the image id 12 in the following `image_data` array: `{..., "image_data": [{"data": "<BASE64_STRING>", "id": 12}]}`. Use `image_data` only with multimodal models, e.g., LLaVA.
	ImageData []ImageData `json:"image_data"`
}

type ImageData struct {
	Data string `json:"data"`
	ID   int    `json:"id"`
}

type GenerationSettings struct {
	CompletionSettings
	Model string `json:"model"`
	NCTX  int    `json:"n_ctx"`
}

type CompletionPostResponse struct {
	Content    string `json:"content"`
	Stop       bool   `json:"stop"`
	Generation string `json:"generation_settings"`
	Model      string `json:"model"`
	Prompt     string `json:"prompt"`
	// StoppedEos indicates whether the completion has stopped because it encountered the EOS token.
	StoppedEos bool `json:"stopped_eos"`
	// StoppedLimit indicates whether the completion stopped because `n_predict` tokens were generated before stop words or EOS was encountered.
	StoppedLimit bool `json:"stopped_limit"`
	// StoppedWord indicates whether the completion stopped due to encountering a stopping word from `stop` JSON array provided.
	StoppedWord bool `json:"stopped_word"`
	// StoppingWord is the stopping word encountered which stopped the generation (or "" if not stopped due to a stopping word).
	StoppingWord string `json:"stopping_word"`
	// Timings is a hash of timing information about the completion such as the number of tokens `predicted_per_second`.
	Timings map[string]float64 `json:"timings"`
	// TokensCached is the number of tokens from the prompt which could be re-used from previous completion (`n_past`).
	TokensCached int `json:"tokens_cached"`
	// TokensEvaluated is the number of tokens evaluated in total from the prompt.
	TokensEvaluated int `json:"tokens_evaluated"`
	// Truncated indicates if the context size was exceeded during generation, i.e. the number of tokens provided in the prompt (`tokens_evaluated`) plus tokens generated (`tokens predicted`) exceeded the context size (`n_ctx`).
	Truncated bool `json:"truncated"`
	// SlotID assigns the completion task to an specific slot. If is -1 the task will be assigned to a Idle slot (default: -1).
	SlotID int `json:"slot_id"`
	// CachePrompt saves the prompt and generation for avoid reprocess entire prompt if a part of this isn't change (default: false).
	CachePrompt bool `json:"cache_prompt"`
	// SystemPrompt changes the system prompt (initial prompt of all slots), this is useful for chat applications. [See more](#change-system-prompt-on-runtime).
	SystemPrompt string `json:"system_prompt"`
}

type EmbeddingPostRequest struct {
	Content string `json:"content"`
}

type EmbeddingPostResponse struct {
	Embedding []float64 `json:"embedding"`
}

type TokenizePostRequest struct {
	Content string `json:"content"`
}

type TokenizePostResponse struct {
	Tokens []int `json:"tokens"`
}

type DetokenizePostRequest struct {
	Tokens []int `json:"tokens"`
}

type DetokenizePostResponse struct {
	Content string `json:"content"`
}
