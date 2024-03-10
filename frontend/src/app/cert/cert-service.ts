import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { CertificateType } from "./cert-data";
import { KeycloakService } from "keycloak-angular";

@Injectable()
export class CertificateService{
    constructor(private http:HttpClient, private kc:KeycloakService){}

    getCerts(){
        return this.http.get<CertificateType[]>("/api/certs", {observe: 'body', responseType: 'json'})
    }
    createCert(){
        var a = this.kc.addTokenToHeader().subscribe();

        return this.http.post("/api/certs/", {})
    }
    getCertByCN(cn:string){
        return this.http.get("/api/certs/" + cn, {responseType: 'blob'});
    }
    revokeByCN(cn:string){
        return this.http.delete("/api/certs/" + cn, {});
    }
}