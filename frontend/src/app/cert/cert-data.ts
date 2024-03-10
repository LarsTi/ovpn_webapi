export type CertificateType = {
    id: number,
    mail: string,
    common_name: string,
    issued_at: Date,
    valid_to: Date
    revoked: boolean
}