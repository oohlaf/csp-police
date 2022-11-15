CREATE TABLE csp_report(
    checksum TEXT PRIMARY KEY,
    document_uri TEXT,
    referrer_uri TEXT,
    violated_directive TEXT,
    effective_directive TEXT,
    original_policy TEXT,
    disposition TEXT,
    blocked_uri TEXT,
    status_code INTEGER,
    source_uri TEXT,
    line_number INTEGER,
    column_number INTEGER,
    script_sample TEXT
);

CREATE TABLE csp_policy(
    checksum TEXT PRIMARY KEY,
    policy TEXT
);

CREATE TABLE csp_log(
    id INTEGER PRIMARY KEY,
    timestamp INTEGER NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT,
    application TEXT,
    policy_version TEXT,
    policy_forced BOOLEAN,
    policy_checksum TEXT,
    report_checksum TEXT NOT NULL,
    FOREIGN KEY (report_checksum) REFERENCES csp_report (report_checksum),
    FOREIGN KEY (policy_checksum) REFERENCES csp_policy (policy_checksum)
);

/*
 TODO

 Normalize _uri as foreign keys to uri table
 uri table can contain host also, like:
  - id
  - uri
  - host (just the domain name without scheme or path)

Normalize user_agent to separate table:
  - id
  - user_agent
  - browser
  - version
  - OS

*/
