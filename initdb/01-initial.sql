CREATE TABLE public.deductions
(
    id serial NOT NULL,
    name character varying(100) NOT NULL,
    amount numeric(10, 2) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT unique_deduction_name UNIQUE (name)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.deductions
    OWNER to postgres;

INSERT INTO public.deductions (name, amount) VALUES ('personal', 60000.0);
INSERT INTO public.deductions (name, amount) VALUES ('donation', 100000.0);
INSERT INTO public.deductions (name, amount) VALUES ('k-receipt', 50000.0);
