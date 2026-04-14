CREATE TABLE reactions (
    user_id UUID NOT NULL,
    chart_id UUID NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (user_id, chart_id),
    CONSTRAINT reactions_type_check CHECK (type IN ('like', 'dislike', 'view'))
);

CREATE INDEX idx_reactions_chart_id ON reactions(chart_id);
CREATE INDEX idx_reactions_user_type ON reactions(user_id, type);