package goproducts.servicestat;

public class Service {
    public final String urlString;
    private ServiceStatus status;
    private final String name;

    public Service(String name, String urlString) {
        this.name = name;
        this.urlString = urlString;
        this.setStatusUnchecked();
    }

    public String getUrlString(){
        return this.urlString;
    }
    public void setStatus(ServiceStatus status) {
        this.status = status;
    }

    public void setStatusActive() {
        this.setStatus(ServiceStatus.ACTIVE);
    }

    public void setStatusInactive() {
        this.setStatus(ServiceStatus.INACTIVE);
    }

    public void setStatusUnchecked() {
        this.setStatus(ServiceStatus.UNCHECKED);
    }

    public ServiceStatus getStatus() {
        return this.status;
    }


    public String getName() {
        return name;
    }
}
